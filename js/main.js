var wordnik = function () {
    var API_ROOT = "http://api.wordnik.com/api/";
    var API_KEY = "36f4cb6b11835b894a0040ebe9209415c1ad48c44764757e2";

    var url = function (path) { return API_ROOT + path + "?api_key=" + API_KEY }

    return {
        getRandom: function (callback) {
            var req = new Request.JSONP({
                url: url("words.json/randomWord") + "&hasDictionaryDef=true",
                onComplete: function (json) {
                    $("word").innerHTML = json.wordstring;
                    if (typeof(callback) === "function") callback(json.wordstring);
                },
            }).send();
        },
        getDefinition: function (word) {
            var req = new Request.JSONP({
                url: url("word.json/" + word + "/definitions"),
                onComplete: function (json) {
                    $("definition").innerHTML = json[0].text;
                },
            }).send();
        }
    }
}

var w = wordnik();
document.addEvent('domready', function(){w.getRandom(w.getDefinition)});

