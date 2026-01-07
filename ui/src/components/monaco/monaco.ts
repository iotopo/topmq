import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker'
import jsonWorker from 'monaco-editor/esm/vs/language/json/json.worker?worker'
// import cssWorker from 'monaco-editor/esm/vs/language/css/css.worker?worker'
// import htmlWorker from 'monaco-editor/esm/vs/language/html/html.worker?worker'
// import tsWorker from 'monaco-editor/esm/vs/language/typescript/ts.worker?worker'
// import sqlWorker from 'monaco-editor/esm/vs/language/sql/sql.worker?worker'

window.MonacoEnvironment = {
  getWorker(_: string, label: string) {
    if (label === 'json') {
      // eslint-disable-next-line new-cap
      return new jsonWorker()
    }
    // if (label === 'typescript' || label === 'javascript') {
    //   return new tsWorker()
    // }

    // eslint-disable-next-line new-cap
    return new editorWorker()
  },
}
