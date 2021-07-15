type SpiderTabName = 'overview' | 'files' | 'tasks' | 'settings';

interface SpiderDialogVisible extends DialogVisible {
  run: boolean;
}

type SpiderDialogKey = DialogKey | 'run';
