import CodeMirror, {Editor, EditorConfiguration} from 'codemirror';

// import addons
import 'codemirror/addon/search/search.js';
import 'codemirror/addon/search/searchcursor';
import 'codemirror/addon/search/matchesonscrollbar.js';
import 'codemirror/addon/search/matchesonscrollbar.css';
import 'codemirror/addon/search/match-highlighter';
import 'codemirror/addon/edit/matchtags';
import 'codemirror/addon/edit/matchbrackets';
import 'codemirror/addon/edit/closebrackets';
import 'codemirror/addon/edit/closetag';
import 'codemirror/addon/comment/comment';
import 'codemirror/addon/hint/show-hint';

// import keymap
import 'codemirror/keymap/emacs.js';
import 'codemirror/keymap/sublime.js';
import 'codemirror/keymap/vim.js';

const themes = [
  '3024-day',
  '3024-night',
  'abcdef',
  'ambiance',
  'ambiance-mobile',
  'ayu-dark',
  'ayu-mirage',
  'base16-dark',
  'base16-light',
  'bespin',
  'blackboard',
  'cobalt',
  'colorforth',
  'darcula',
  'dracula',
  'duotone-dark',
  'duotone-light',
  'eclipse',
  'elegant',
  'erlang-dark',
  'gruvbox-dark',
  'hopscotch',
  'icecoder',
  'idea',
  'isotope',
  'lesser-dark',
  'liquibyte',
  'lucario',
  'material',
  'material-darker',
  'material-ocean',
  'material-palenight',
  'mbo',
  'mdn-like',
  'midnight',
  'monokai',
  'moxer',
  'neat',
  'neo',
  'night',
  'nord',
  'oceanic-next',
  'panda-syntax',
  'paraiso-dark',
  'paraiso-light',
  'pastel-on-dark',
  'railscasts',
  'rubyblue',
  'seti',
  'shadowfox',
  'solarized',
  'ssms',
  'the-matrix',
  'tomorrow-night-bright',
  'tomorrow-night-eighties',
  'ttcn',
  'twilight',
  'vibrant-ink',
  'xq-dark',
  'xq-light',
  'yeti',
  'yonce',
  'zenburn',
];

const template = `import os
def func(a):
  pass
class Class1:
  pass
`;

const optionsDefinitions: FileEditorOptionDefinition[] = [
  {
    name: 'theme',
    title: 'Theme',
    description: 'The theme to style the editor with.',
    type: 'select',
    data: {
      options: themes,
    },
  },
  {
    name: 'indentUnit',
    title: 'Indent Unit',
    description: 'How many spaces a block (whatever that means in the edited language) should be indented.',
    type: 'input-number',
    data: {
      min: 1,
    }
  },
  {
    name: 'smartIndent',
    title: 'Smart Indent',
    description: 'Whether to use the context-sensitive indentation that the mode provides (or just indent the same as the line before).',
    type: 'switch',
  },
  {
    name: 'tabSize',
    title: 'Tab Size',
    description: 'The width of a tab character. Defaults to 4.',
    type: 'input-number',
    data: {
      min: 1,
    }
  },
  {
    name: 'indentWithTabs',
    title: 'Indent with Tabs',
    description: 'Whether, when indenting, the first N*tabSize spaces should be replaced by N tabs.',
    type: 'switch',
  },
  {
    name: 'electricChars',
    title: 'Electric Chars',
    description: 'Configures whether the editor should re-indent the current line when a character is typed that might change its proper indentation (only works if the mode supports indentation).',
    type: 'switch',
  },
  {
    name: 'keyMap',
    title: 'Keymap',
    description: 'Configures the keymap to use.',
    type: 'select',
    data: {
      options: [
        'default',
        'emacs',
        'sublime',
        'vim',
      ]
    },
  },
  {
    name: 'lineWrapping',
    title: 'Line Wrapping',
    description: 'Whether to scroll or wrap for long lines.',
    type: 'switch',
  },
  {
    name: 'lineNumbers',
    title: 'Line Numbers',
    description: 'Whether to show line numbers to the left of the editor.',
    type: 'switch',
  },
  {
    name: 'showCursorWhenSelecting',
    title: 'Show Cursor When Selecting',
    description: 'Whether the cursor should be drawn when a selection is active.',
    type: 'switch',
  },
  {
    name: 'lineWiseCopyCut',
    title: 'Line-wise Copy-Cut',
    description: 'When enabled, doing copy or cut when there is no selection will copy or cut the whole lines that have cursors on them.',
    type: 'switch',
  },
  {
    name: 'pasteLinesPerSelection',
    title: 'Paste Lines per Selection',
    description: 'When pasting something from an external source (not from the editor itself), if the number of lines matches the number of selection, the editor will by default insert one line per selection. You can set this to false to disable that behavior.',
    type: 'switch',
  },
  {
    name: 'undoDepth',
    title: 'Paste Lines per Selection',
    description: 'The maximum number of undo levels that the editor stores.',
    type: 'input-number',
    data: {
      min: 1,
    },
  },
  {
    name: 'cursorBlinkRate',
    title: 'Cursor Blink Rate',
    description: 'Half-period in milliseconds used for cursor blinking.',
    type: 'input-number',
    data: {
      min: 10,
    },
  },
  {
    name: 'cursorScrollMargin',
    title: 'Cursor Scroll Margin',
    description: 'How much extra space to always keep above and below the cursor when approaching the top or bottom of the visible view in a scrollable document.',
    type: 'input-number',
    data: {
      min: 0,
    },
  },
  {
    name: 'cursorHeight',
    title: 'Cursor Height',
    description: 'Determines the height of the cursor. Setting it to 1, means it spans the whole height of the line. For some fonts (and by some tastes) a smaller height (for example 0.85), which causes the cursor to not reach all the way to the bottom of the line, looks better',
    type: 'input-number',
    data: {
      min: 0,
      step: 0.01,
    },
  },
  {
    name: 'maxHighlightLength',
    title: 'Max Highlight Length',
    description: 'When highlighting long lines, in order to stay responsive, the editor will give up and simply style the rest of the line as plain text when it reaches a certain position.',
    type: 'input-number',
    data: {
      min: 1,
    },
  },
  {
    name: 'spellcheck',
    title: 'Spell Check',
    description: 'Specifies whether or not spellcheck will be enabled on the input.',
    type: 'switch',
  },
  {
    name: 'autocorrect',
    title: 'Auto Correct',
    description: 'Specifies whether or not auto-correct will be enabled on the input.',
    type: 'switch',
  },
  {
    name: 'autocapitalize',
    title: 'Auto Capitalize',
    description: 'Specifies whether or not auto-capitalization will be enabled on the input.',
    type: 'switch',
  },
  {
    name: 'highlightSelectionMatches',
    title: 'Highlight Selection Matches',
    description: 'Adds a highlightSelectionMatches option that can be enabled to highlight all instances of a currently selected word. When enabled, it causes the current word to be highlighted when nothing is selected.',
    type: 'switch',
  },
  {
    name: 'matchBrackets',
    title: 'Match Brackets',
    description: 'When set to true or an options object, causes matching brackets to be highlighted whenever the cursor is next to them.',
    type: 'switch',
  },
  {
    name: 'matchTags',
    title: 'Match Tags',
    description: 'When enabled will cause the tags around the cursor to be highlighted',
    type: 'switch',
  },
  {
    name: 'autoCloseBrackets',
    title: 'Auto-Close Brackets',
    description: 'Will auto-close brackets and quotes when typed. It\'ll auto-close ()[]{}\'\'"".',
    type: 'switch',
  },
  {
    name: 'autoCloseTags',
    title: 'Auto-Close Tags',
    description: 'Will auto-close XML tags when \'>\' or \'/\' is typed.',
    type: 'switch',
  },
  {
    name: 'showHint',
    title: 'Show Hint',
    description: '',
    type: 'switch',
  },
];

const themeCache = new Set<string>();

export const getCodemirrorEditor = (el: HTMLElement, options: EditorConfiguration): Editor => {
  return CodeMirror(el, options);
};

export const initTheme = async (name?: string) => {
  if (!name) name = 'darcula';
  if (themeCache.has(name)) return;
  await import(`codemirror/theme/${name}.css`);
  themeCache.add(name);
};

export const getThemes = () => {
  return themes;
};

export const getCodeMirrorTemplate = () => {
  return template;
};

export const getOptionDefinition = (name: string): FileEditorOptionDefinition | undefined => {
  return optionsDefinitions.find(d => d.name === name);
};
