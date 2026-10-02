'use strict';

const { cleanHTML } = require('./src/clean-html');

module.exports = { cleanHTML, default: cleanHTML };

// When executed directly, act as a CLI mirroring the Go --clean-html flag.
//   node index.js --clean-html "<p>hi</p>"
//   node index.js --clean-html -            (read stdin)
//   cat page.html | node index.js
if (require.main === module) {
  const argv = process.argv.slice(2);
  const getFlag = (name) => {
    const eq = argv.find((a) => a.startsWith(name + '='));
    if (eq) return eq.slice(name.length + 1);
    const i = argv.indexOf(name);
    return i >= 0 ? argv[i + 1] : undefined;
  };

  const htmlArg = getFlag('--clean-html');
  const fileArg = getFlag('--clean-html-file');

  const run = () => {
    if (fileArg) {
      const fs = require('fs');
      process.stdout.write(cleanHTML(fs.readFileSync(fileArg, 'utf8')));
      return;
    }
    if (htmlArg === '-') {
      const data = require('fs').readFileSync(0, 'utf8');
      process.stdout.write(cleanHTML(data));
      return;
    }
    if (htmlArg) {
      process.stdout.write(cleanHTML(htmlArg));
      return;
    }
    // Fall back to piped stdin.
    if (!process.stdin.isTTY) {
      const data = require('fs').readFileSync(0, 'utf8');
      process.stdout.write(cleanHTML(data));
      return;
    }
    process.stderr.write(
      'Usage: clean-html --clean-html "<html>" | --clean-html-file <path> | cat file | clean-html\n'
    );
    process.exitCode = 1;
  };

  run();
}
