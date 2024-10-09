package notification

const defaultTheme = `<style>
.editor-container {
  color: #000;
  position: relative;
  font-weight: 400;
  text-align: left;
  height: 100%;
}

.editor-inner {
  height: calc(100% - 45px);
  background: #fff;
  position: relative;
}

.editor-placeholder {
  color: #999;
  overflow: hidden;
  position: absolute;
  text-overflow: ellipsis;
  top: 15px;
  left: 10px;
  font-size: 15px;
  user-select: none;
  display: inline-block;
  pointer-events: none;
}

.editor-input {
  height: 100%;
  resize: none;
  font-size: 15px;
  position: relative;
  tab-size: 1;
  outline: 0;
  padding: 15px 10px;
  caret-color: #444;
}

.ltr {
  text-align: left;
}

.rtl {
  text-align: right;
}

.editor-text-bold {
  font-weight: bold;
}

.editor-text-italic {
  font-style: italic;
}

.editor-text-underline {
  text-decoration: underline;
}

.editor-text-strikethrough {
  text-decoration: line-through;
}

.editor-text-underlineStrikethrough {
  text-decoration: underline line-through;
}

.editor-text-code {
  background-color: rgb(240, 242, 245);
  padding: 1px 0.25rem;
  font-family: Menlo, Consolas, Monaco, monospace;
  font-size: 94%;
}

.editor-link {
  color: rgb(33, 111, 219);
  text-decoration: none;

  &:hover {
    color: rgb(33, 111, 219);
    text-decoration: underline;
  }
}

.tree-view-output {
  display: block;
  background: #222;
  color: #fff;
  padding: 5px;
  font-size: 12px;
  white-space: pre-wrap;
  margin: 1px auto 10px auto;
  max-height: 250px;
  position: relative;
  border-bottom-left-radius: 10px;
  border-bottom-right-radius: 10px;
  overflow: auto;
  line-height: 14px;
}

.editor-code {
  background-color: rgb(240, 242, 245);
  font-family: Menlo, Consolas, Monaco, monospace;
  display: block;
  padding: 8px 8px 8px 52px;
  line-height: 1.53;
  font-size: 13px;
  margin: 8px 0;
  tab-size: 2;
  overflow-x: auto;
  position: relative;
}

.editor-code:before {
  content: attr(data-gutter);
  position: absolute;
  background-color: #eee;
  left: 0;
  top: 0;
  border-right: 1px solid #ccc;
  padding: 8px;
  color: #777;
  white-space: pre-wrap;
  text-align: right;
  min-width: 25px;
}

.editor-code:after {
  content: attr(data-highlight-language);
  top: 0;
  right: 3px;
  padding: 3px;
  font-size: 10px;
  text-transform: uppercase;
  position: absolute;
  color: rgba(0, 0, 0, 0.5);
}

.editor-tokenComment {
  color: slategray;
}

.editor-tokenPunctuation {
  color: #999;
}

.editor-tokenProperty {
  color: #905;
}

.editor-tokenSelector {
  color: #690;
}

.editor-tokenOperator {
  color: #9a6e3a;
}

.editor-tokenAttr {
  color: #07a;
}

.editor-tokenVariable {
  color: #e90;
}

.editor-tokenFunction {
  color: #dd4a68;
}

.editor-paragraph {
  margin: 0 0 8px;
  position: relative;
}

.editor-paragraph:last-child {
  margin-bottom: 0;
}

.editor-heading-h1 {
  font-size: 24px;
  color: rgb(5, 5, 5);
  font-weight: 400;
  margin: 0 0 12px;
  padding: 0;
}

.editor-heading-h2 {
  font-size: 20px;
  color: rgb(101, 103, 107);
  font-weight: 700;
  margin: 10px 0 0;
  padding: 0;
}

.editor-heading-h3 {
  font-size: 16px;
  color: rgb(101, 103, 107);
  font-weight: 600;
  margin: 10px 0 0;
  padding: 0;
}

.editor-heading-h4 {
  font-size: 14px;
  color: rgb(101, 103, 107);
  font-weight: 500;
  margin: 10px 0 0;
  padding: 0;
}

.editor-quote {
  margin: 0 0 0 20px;
  font-size: 15px;
  color: rgb(101, 103, 107);
  border-left-color: rgb(206, 208, 212);
  border-left-width: 4px;
  border-left-style: solid;
  padding-left: 16px;
}

.editor-list-ol {
  padding: 0;
  margin: 0 0 0 16px;
}

.editor-list-ul {
  padding: 0;
  margin: 0 0 0 16px;
}

.editor-listitem {
  margin: 8px 32px 8px 32px;
}

.editor-nested-listitem {
  list-style-type: none;
}

.editor-table {
  border-collapse: collapse;
  border-spacing: 0;
  overflow-y: scroll;
  overflow-x: scroll;
  table-layout: auto;
  width: max-content;
  margin: 0 25px 30px 0;
}

.editor-cell {
  position: relative;
  border: 1px solid #bbb;
  min-width: 150px;
  vertical-align: top;
  text-align: start;
  padding: 6px 8px;
  outline: none;
}
</style>`

func GetTheme() string {
	return defaultTheme
}
