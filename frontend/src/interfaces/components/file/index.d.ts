import {EditorConfiguration} from 'codemirror';
import '@/utils/codemirror';

declare global {
  interface FileNavItem {
    id?: string;
    is_dir?: boolean;
    path?: string;
    name?: string;
    extension?: string;
    children?: FileNavItem[];
  }

  interface FileEditorStyle {
    backgroundColor?: string;
    color?: string;
    height?: string;
    backgroundColorGutters?: string;
  }

  interface FileEditorConfigurationSearch {
    bottom?: boolean;
  }

  interface FileEditorConfiguration extends EditorConfiguration {
    search?: FileEditorConfigurationSearch;
  }

  type FileEditorOptionDefinitionType = 'select' | 'input-number' | 'switch';

  interface FileEditorOptionDefinitionData {
    options?: string[];
    min?: number;
    step?: number;
  }

  interface FileEditorOptionDefinition {
    name: string;
    title: string;
    description: string;
    type: FileEditorOptionDefinitionType;
    data?: FileEditorOptionDefinitionData;
  }
}
