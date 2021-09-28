/*
  Processes results while inside an editing query page in Stack Exchange Data Explorer
  It must be executed under the console pane and user must also be logged in
  Permission needs to be given(allowed) to the page for downloading many files
*/

let tagNames = [
  ["rx-java", "rx-java2", "rx-java3"], //rxjava
  ["rxjs", "rxjs5", "rxjs6","rxjs7"], //rxjs
  ["rx-kotlin", "rx-kotlin2"], //rxkotlin
  ["rx-swift"] //rxswift
];

// loads the in-page CodeMirror instance
const editor = document.querySelector('.CodeMirror').CodeMirror;

function buildQuery(tagName) {
  return `select p.* from posts p inner join PostTags ps
on p.Id = ps.PostId inner join Tags t on ps.TagId = t.Id and
t.TagName = '${tagName}' where p.ParentId is null order by p.Id`;
}

function writeQuery(query) {
  editor.clearHistory();
  editor.setValue(query);
}

function timeoutPromiseResolve(interval) {
  return new Promise((resolve, reject) => {
    setTimeout(function () {
      resolve();
    }, interval);
  });
};

function verifyResult(timer) {
  return new Promise((resolve, reject) => {
    const interval = setInterval(function () {
      let errorElem = document.getElementById("error-message");
      let loading = document.getElementById("loading");
      if (errorElem.style.display != 'none') {
        reject(new Error(errorElem));
        clearInterval(interval);
      } else {
        if (loading.style.display == 'none') {
          resolve();
          clearInterval(interval);
        }
      }
    }, timer);
  });
};

async function processQuery(query) {
  console.log("Processing query...");
  writeQuery(query);
  document.getElementById("submit-query").click();
  await verifyResult(5000); //5s
}

async function executeQuery(dist = 0) {
  try {
    for (let j = 0; j < tagNames[dist].length; j++) {
      await processQuery(buildQuery(tagNames[dist][j]));
      document.getElementById("resultSetsButton").click();
      await timeoutPromiseResolve(5000); //5s
    }
    console.log("Done!");
  } catch (e) {
    console.log(e);
  }
}

executeQuery(0); //rxjava
//executeQuery(1); //rxjs
//executeQuery(2); //rxkotlin
//executeQuery(3); //rxswift