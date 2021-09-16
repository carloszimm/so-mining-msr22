/*
	Processes results while inside an editing query page in Stack Exchange Data Explorer
	Permission needs to be given(allowed) to the page for downloading many files
*/

let tagNames = [
["rxjs","rxjs5","rxjs6","rxjs-observables","rxjs-pipeable-operators",
"reactive-extensions-js","rxjs-marbles","rxjs-dom","rxjs-lettable-operators","rxjs-subscriptions",
"rxfire","rxjs-compat","rxjs-test-scheduler","rxjs-pipe","rxjs7","rx-angular"],
["rx-swift"]];

// loads the in-page CodeMirror instance
let editor = document.querySelector('.CodeMirror').CodeMirror;

function buildQuery(tagName){
return `select p.* from posts p inner join PostTags ps
on p.Id = ps.PostId inner join Tags t on ps.TagId = t.Id and
t.TagName = '${tagName}' where p.ParentId is null order by p.Id`;
}

function writeQuery(query){
	editor.clearHistory();
	editor.setValue(query);
}

function timeoutPromiseResolve(interval) {
  return new Promise((resolve, reject) => {
    setTimeout(function(){
      resolve();
    }, interval);
  });
};

function verifyResult(timer) {
  return new Promise((resolve, reject) => {
    const interval = setInterval(function(){
		let errorElem = document.getElementById("error-message");
		let loading = document.getElementById("loading");
		if(errorElem.style.display != 'none'){
			reject(new Error(errorElem));
			clearInterval(interval);
		}else {
			if(loading.style.display == 'none'){
				resolve();
				clearInterval(interval);
			}
		}
    }, timer);
  });
};

async function processQuery(query){
	console.log("Processing query...");
	writeQuery(query);
	document.getElementById("submit-query").click();
	await verifyResult(5000); //5s
}
  
async function executeQuery(){
	try{
		for(let i = 0; i<tagNames.length; i++){
			for(let j = 0; j<tagNames[i].length; j++){
				await processQuery(buildQuery(tagNames[i][j]));
				document.getElementById("resultSetsButton").click();
				await timeoutPromiseResolve(10000); //10s
			}
		}
	}catch(e){
		console.log(e);
	}
}

executeQuery();