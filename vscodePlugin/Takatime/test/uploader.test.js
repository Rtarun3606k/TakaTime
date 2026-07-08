const assert = require('assert');

// Mock vscode module
const mockVscode = {
  env: {
    appName: 'Visual Studio Code'
  },
  workspace: {
    getWorkspaceFolder: () => null
  }
};

// Override require for vscode
const Module = require('module');
const originalRequire = Module.prototype.require;
Module.prototype.require = function(id) {
  if (id === 'vscode') {
    return mockVscode;
  }
  return originalRequire.apply(this, arguments);
};

// Now require the uploader module
const { getGoArgs } = require('../Plugin/Uploader');

suite('Uploader Test Suite', () => {
  test('getGoArgs returns VS Code editor flag for standard VS Code', () => {
    mockVscode.env.appName = 'Visual Studio Code';
    const mockDoc = { fileName: '/test/file.js', languageId: 'javascript' };
    const args = getGoArgs(mockDoc, 'mongodb://test');
    const editorIndex = args.indexOf('-editor');
    assert.notStrictEqual(editorIndex, -1);
    assert.strictEqual(args[editorIndex + 1], 'VS Code');
  });

  test('getGoArgs returns VS Code editor flag for VS Code Insiders', () => {
    mockVscode.env.appName = 'Visual Studio Code - Insiders';
    const mockDoc = { fileName: '/test/file.js', languageId: 'javascript' };
    const args = getGoArgs(mockDoc, 'mongodb://test');
    const editorIndex = args.indexOf('-editor');
    assert.notStrictEqual(editorIndex, -1);
    assert.strictEqual(args[editorIndex + 1], 'VS Code');
  });

  test('getGoArgs returns Cursor editor flag for Cursor', () => {
    mockVscode.env.appName = 'Cursor';
    const mockDoc = { fileName: '/test/file.js', languageId: 'javascript' };
    const args = getGoArgs(mockDoc, 'mongodb://test');
    const editorIndex = args.indexOf('-editor');
    assert.notStrictEqual(editorIndex, -1);
    assert.strictEqual(args[editorIndex + 1], 'Cursor');
  });

  test('getGoArgs returns Windsurf editor flag for Windsurf', () => {
    mockVscode.env.appName = 'Windsurf';
    const mockDoc = { fileName: '/test/file.js', languageId: 'javascript' };
    const args = getGoArgs(mockDoc, 'mongodb://test');
    const editorIndex = args.indexOf('-editor');
    assert.notStrictEqual(editorIndex, -1);
    assert.strictEqual(args[editorIndex + 1], 'Windsurf');
  });

  test('getGoArgs returns Antigravity for unknown editor', () => {
    mockVscode.env.appName = 'SomeUnknownEditor';
    const mockDoc = { fileName: '/test/file.js', languageId: 'javascript' };
    const args = getGoArgs(mockDoc, 'mongodb://test');
    const editorIndex = args.indexOf('-editor');
    assert.notStrictEqual(editorIndex, -1);
    assert.strictEqual(args[editorIndex + 1], 'Antigravity');
  });
});
