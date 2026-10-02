'use strict';

// Test runner for the JavaScript clean-html port. It validates against the
// same canonical vectors used by the Go and Python implementations, so a
// passing run here means all three ports agree on behaviour.

const assert = require('assert');
const fs = require('fs');
const path = require('path');
const { cleanHTML } = require('../src/clean-html');

const vectors = JSON.parse(
  fs.readFileSync(path.join(__dirname, '..', '..', 'testdata', 'clean-html-cases.json'), 'utf8')
);

const trimLastNewline = (s) => (s.endsWith('\n') ? s.slice(0, -1) : s);

let failed = 0;
for (const c of vectors.cases) {
  const got = trimLastNewline(cleanHTML(c.in));
  try {
    assert.strictEqual(got, c.want);
    process.stdout.write(`ok   - ${c.name}\n`);
  } catch (err) {
    failed += 1;
    process.stdout.write(`FAIL - ${c.name}\n  input:  ${JSON.stringify(c.in)}\n  got:    ${JSON.stringify(got)}\n  want:   ${JSON.stringify(c.want)}\n`);
  }
}

if (failed > 0) {
  process.stdout.write(`\n${failed} case(s) failed\n`);
  process.exit(1);
}
process.stdout.write(`\nAll ${vectors.cases.length} cases passed\n`);
