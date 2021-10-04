const fs = require("fs"),
	path = require("path");

const dist = parseInt(process.argv[2]) || 0;

let tagNames = [
	["rx-java", "rx-java2", "rx-java3"], //rxjava
	["rxjs", "rxjs5", "rxjs6", "rxjs7"], //rxjs
	["rx-kotlin", "rx-kotlin2"], //rxkotlin
	["rx-swift"] //rxswift
];

const srcFolder = path.resolve(__dirname, "./staging area");
const destFolder = path.resolve(__dirname, `../../assets/data explorer/${getDistName()}`)

function getDistName() {
	switch (dist) {
		case 1:
			return "rxjs";
		case 2:
			return "rxkotlin";
		case 3:
			return "rxswift";
		default:
			return "rxjava";
	}
}

fs.rmdir(destFolder, { recursive: true, force: true }, err => {
	if (err) {
		console.error(err);
		process.exit(1);
	}

	fs.mkdirSync(destFolder);

	fs.readdir(srcFolder, (err, files) => {
		if (err) {
			throw err;
			process.exit(1);
		}
		files.forEach(file => {
			let fileName = path.basename(file, ".csv");

			if (fileName.includes("QueryResults")) {
				let index = fileName.replace(/[^0-9]+/g, '');

				index = index === "" ? 0 : parseInt(index);

				if (index < tagNames[dist].length) {
					rename(file, "posts_" + tagNames[dist][index] + ".csv")
				} else {
					let i = index % tagNames[dist].length;
					rename(file, "posts_" + tagNames[dist][i] + "_withAnswers.csv")
				}
			}
		});
	});
});

function rename(file, name) {
	fs.rename(path.join(srcFolder, file),
		path.join(destFolder, name),
		err => {
			if (err) {
				throw err;
				process.exit(1);
			}
		});
}