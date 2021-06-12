const fs = require('fs');

let res1 = 0;
let res2 = 0;
let rawdata = fs.readFileSync('./day12/input');
let data = JSON.parse(rawdata);
parseObj(data, false);
console.log("Part 1: ", res1);
console.log("Part 2: ", res2);

function parseObj(json, b) {
    let keys = Object.keys(json);
    if (Object.values(json).includes("red")) b = true;
    keys.forEach(elem => {
        let e = json[elem]
        if (!isNaN(e)) {
            res1 += parseInt(e);
            if (!b) res2 += parseInt(e);
        }
        else if (typeof e == "object") {
            if (Array.isArray(e)) {
                parseArray(e, b);
            }
            else {
                parseObj(e, b);
            }
        }
    });
};

function parseArray(arr, b) {
    arr.forEach(e => {
        if (!isNaN(e)) {
            res1 += parseInt(e);
            if (!b) res2 += parseInt(e);
        }
        else if (typeof e == "object") {
            if (Array.isArray(e)) {
                parseArray(e, b);
            }
            else {
                parseObj(e, b);
            }
        }
    });
};