'use strict';

// clean-html.js — HTML to Markdown cleaner.
//
// This is a behaviour-identical port of the Go implementation in
// internal/markdown/markdown.go. All three language ports (Go, JavaScript,
// Python) are validated against the same canonical vectors in
// ../../testdata/clean-html-cases.json.

// Base regexes are defined WITHOUT the global flag so .exec is stateless.
// Helper replaceAll() adds `g` when replacing every occurrence.
const RE = {
  comment: /<!--[\s\S]*?-->/i,
  script: /<script\b[^>]*>[\s\S]*?<\/script\s*>/i,
  style: /<style\b[^>]*>[\s\S]*?<\/style\s*>/i,
  alt: /\balt\s*=\s*["']([^"']*)["']/i,

  heading: /<h([1-6])\b[^>]*>([\s\S]*?)<\/h[1-6]\s*>/i,
  image: /<img\b[^>]*?src\s*=\s*["']([^"']*)["'][^>]*?>/i,
  link: /<a\b[^>]*?href\s*=\s*["']([^"']*)["'][^>]*>([\s\S]*?)<\/a\s*>/i,

  bold: /<(?:strong|b)\b[^>]*>([\s\S]*?)<\/(?:strong|b)\s*>/i,
  italic: /<(?:em|i)\b[^>]*>([\s\S]*?)<\/(?:em|i)\s*>/i,
  code: /<code\b[^>]*>([\s\S]*?)<\/code\s*>/i,
  pre: /<pre\b[^>]*>([\s\S]*?)<\/pre\s*>/i,

  quote: /<blockquote\b[^>]*>([\s\S]*?)<\/blockquote\s*>/i,
  list: /<(ul|ol)\b[^>]*>([\s\S]*?)<\/(?:ul|ol)\s*>/i,
  item: /<li\b[^>]*>([\s\S]*?)<\/li\s*>/i,

  br: /<br\b[^>]*\/?>/i,
  hr: /<hr\b[^>]*\/?>/i,

  block: /<\/?(?:p|div|section|article|header|footer|main|aside|nav|figure|figcaption|table|thead|tbody|tr|td|th|dl|dt|dd|form|fieldset)\b[^>]*>/i,
  anyTag: /<\/?[a-zA-Z][^>]*>/i,

  multiBlank: /\n{3,}/,
  trailingSP: /[ \t]+$/m,
  multiSpaces: /[ \t]{2,}/,
};

function globalRe(re) {
  return new RegExp(re.source, re.flags.indexOf('g') >= 0 ? re.flags : re.flags + 'g');
}

// Replace every match using a callback (like Go's ReplaceAllStringFunc).
function replaceAllFunc(str, re, fn) {
  return str.replace(globalRe(re), fn);
}

// Replace every match with a literal/replacement string.
function replaceAll(str, re, withStr) {
  return str.replace(globalRe(re), withStr);
}

// Minimal HTML entity decoder (named subset + numeric), matching the
// behaviour of Go's html.UnescapeString for the entities we rely on.
const NAMED = { amp: '&', lt: '<', gt: '>', quot: '"', apos: "'", nbsp: '\u00a0' };

function decodeEntities(s) {
  return s.replace(/&(#x?[0-9a-fA-F]+|[a-zA-Z]+);/g, (m, ent) => {
    if (ent[0] === '#') {
      let cp;
      if (ent[1] === 'x' || ent[1] === 'X') cp = parseInt(ent.slice(2), 16);
      else cp = parseInt(ent.slice(1), 10);
      if (Number.isFinite(cp) && cp >= 0 && cp <= 0x10ffff) {
        try {
          return String.fromCodePoint(cp);
        } catch (_) {
          return m;
        }
      }
      return m;
    }
    const k = ent.toLowerCase();
    return Object.prototype.hasOwnProperty.call(NAMED, k) ? NAMED[k] : m;
  });
}

// preStore holds fenced code blocks extracted before generic tag stripping.
let preStore = [];

function convertPre(s) {
  return replaceAllFunc(s, RE.pre, (m) => {
    const inner = RE.pre.exec(m)[1];
    const code = codeText(inner).replace(/\n+$/, '');
    const fence = '```\n' + code + '\n```';
    preStore.push(fence);
    return '\n\x00PRE' + (preStore.length - 1) + '\x00\n';
  });
}

function restorePre(s) {
  return replaceAllFunc(s, /\x00PRE(\d+)\x00/, (tok) => {
    const i = parseInt(/\d+/.exec(tok)[0], 10);
    return i >= 0 && i < preStore.length ? preStore[i] : '';
  });
}

function convertHeadings(s) {
  return replaceAllFunc(s, RE.heading, (m) => {
    const sub = RE.heading.exec(m);
    let level = parseInt(sub[1], 10);
    if (!(level >= 1 && level <= 6)) level = 1;
    const text = inlineText(sub[2]);
    if (!text) return '';
    return '\n' + '#'.repeat(level) + ' ' + text + '\n';
  });
}

function convertImages(s) {
  return replaceAllFunc(s, RE.image, (m) => {
    const src = RE.image.exec(m)[1].trim();
    const a = RE.alt.exec(m);
    const alt = a ? a[1].trim() : '';
    if (!src) return '';
    return '![' + alt + '](' + src + ')';
  });
}

function convertLinks(s) {
  return replaceAllFunc(s, RE.link, (m) => {
    const sub = RE.link.exec(m);
    const href = sub[1].trim();
    const text = inlineText(sub[2]);
    if (!href) return text;
    if (!text) return href;
    if (text === href) return href;
    return '[' + text + '](' + href + ')';
  });
}

function convertLists(s) {
  return replaceAllFunc(s, RE.list, (m) => {
    const sub = RE.list.exec(m);
    const ordered = sub[1].toLowerCase() === 'ol';
    let idx = 0;
    const out = replaceAllFunc(s2(sub[2]), RE.item, (item) => {
      const text = inlineText(RE.item.exec(item)[1]);
      if (!text) return '';
      if (ordered) {
        idx += 1;
        return '\n' + idx + '. ' + text;
      }
      return '\n- ' + text;
    });
    if (!out) return '';
    return '\n' + out.replace(/^\n+/, '') + '\n';
  });
}

function convertQuotes(s) {
  return replaceAllFunc(s, RE.quote, (m) => {
    const text = inlineText(RE.quote.exec(m)[1]);
    if (!text) return '';
    return '\n' + text.split('\n').map((l) => '> ' + l).join('\n') + '\n';
  });
}

function convertInline(s) {
  s = replaceAllFunc(s, RE.bold, (m) => {
    const text = plainText(RE.bold.exec(m)[1]);
    return text ? '**' + text + '**' : '';
  });
  s = replaceAllFunc(s, RE.italic, (m) => {
    const text = plainText(RE.italic.exec(m)[1]);
    return text ? '_' + text + '_' : '';
  });
  s = replaceAllFunc(s, RE.code, (m) => {
    const text = plainText(RE.code.exec(m)[1]);
    return text ? '`' + text + '`' : '';
  });
  return s;
}

function inlineText(s) {
  s = replaceAll(s, RE.bold, '$1');
  s = replaceAll(s, RE.italic, '$1');
  s = replaceAll(s, RE.code, '$1');
  s = replaceAll(s, RE.link, '$2');
  s = replaceAll(s, RE.image, '');
  s = replaceAll(s, RE.br, ' ');
  s = replaceAll(s, RE.anyTag, '');
  s = decodeEntities(s);
  return s.replace(RE.multiSpaces, ' ').trim();
}

function plainText(s) {
  s = replaceAll(s, RE.anyTag, '');
  s = decodeEntities(s);
  return s.replace(RE.multiSpaces, ' ').trim();
}

function codeText(s) {
  s = replaceAll(s, RE.br, '\n');
  s = replaceAll(s, RE.anyTag, '');
  return decodeEntities(s);
}

function normalise(s) {
  s = restorePre(s);
  s = s.split('\u00a0').join(' ');
  s = replaceAll(s, RE.trailingSP, '');
  s = replaceAll(s, RE.multiBlank, '\n\n');
  return s.trim() + '\n';
}

// Small helper kept for readability (identity); mirrors Go passing inner text.
function s2(x) {
  return x;
}

/**
 * Convert an HTML string into clean Markdown.
 * @param {string} input HTML source.
 * @returns {string} Markdown with a trailing newline.
 */
function cleanHTML(input) {
  if (typeof input !== 'string' || input.trim() === '') return '';
  preStore = [];
  let s = input.split('\r\n').join('\n').split('\r').join('\n');
  s = replaceAll(s, RE.comment, '');
  s = replaceAll(s, RE.script, '');
  s = replaceAll(s, RE.style, '');

  s = convertPre(s);
  s = convertHeadings(s);
  s = convertImages(s);
  s = convertLinks(s);
  s = convertLists(s);
  s = convertQuotes(s);
  s = convertInline(s);

  s = replaceAll(s, RE.br, '\n');
  s = replaceAll(s, RE.hr, '\n---\n');
  s = replaceAll(s, RE.block, '\n\n');

  s = replaceAll(s, RE.anyTag, '');
  s = decodeEntities(s);

  return normalise(s);
}

module.exports = { cleanHTML };
