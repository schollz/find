
var user_info;
var message = "Framework for Internal Navigation and Discovery, Hypercube Platforms";
var id = null;
var learnCounter = 0;
var learnLimit = 100;
var progress = 0.0;
var groupfind;
var locationfind;
var userfind;
var currentLocation = "unknown";

function toTitleCase(str)
{
    return str.replace(/\w\S*/g, function(txt){return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase();});
}

function makeAnimalDoublet() {
    adjectives = ["lulling", "vile", "foul", "cheerful", "messy", "dreadful", "uneven", "stinky", "young", "sparkling", "sweltering", "verdant", "hideous", "friendly", "blistering", "rambunctious", "carefree", "fat", "sloppy", "gloomy", "awful", "anemic", "minute", "stiff", "benevolent", "ceaseless", "large", "quick", "round", "glassy", "rusty", "scarce", "odd", "shining", "even", "dowdy", "solemn", "scorching", "brief", "rotten", "new", "plush", "cozy", "meandering", "apologetic", "nimble", "busy", "strong", "great", "brilliant", "piercing", "creepy", "miniature", "narrow", "whimsical", "fantastic", "cowardly", "disgusting", "marvelous", "snug", "stern", "stingy", "angry", "spiky", "cheeky", "gorgeous", "mysterious", "flat", "clever", "charming", "dismal", "meek", "somber", "sour", "thin", "beautiful", "stubborn", "crazy", "challenging", "gaunt", "salty", "indifferent", "huge", "daring", "awkward", "picturesque", "copious", "glowing", "truthful", "rude", "petite", "cranky", "ornery", "brazen", "modest", "purring", "filthy", "rotund", "short", "splendid", "hasty", "deafening", "crawling", "furious", "tall", "mute", "ghastly", "still", "arrogant", "crabby", "haughty", "curly", "voiceless", "hot", "courageous", "late", "microscopic", "vast", "stifling", "good", "disrespectful", "bashful", "moaning", "towering", "adventurous", "idyllic", "careless", "shocking", "erratic", "heavy", "square", "hard", "jagged", "gullible", "hushed", "revolting", "content", "thrifty", "sluggish", "eternal", "greedy", "pleasant", "small", "demanding", "greasy", "enormous", "hostile", "terrible", "chilly", "speedy", "massive", "loud", "puny", "striking", "clumsy", "soaring", "brave", "fast", "delicious", "effortless", "bland", "thick", "little", "ancient", "silent", "wonderful", "stuffed", "grimy", "bitter", "muggy", "shallow", "ridiculous", "absentminded", "fuzzy", "peculiar", "mangy", "wide", "kind", "squeaky", "screeching", "silly", "squealing", "spoiled", "gigantic", "happy", "steep", "ingenious", "modern", "juicy", "gentle", "medium", "brawny", "curved", "lumpy", "afraid", "amusing", "thundering", "cooing", "oppressive", "swollen", "grave", "sturdy", "average", "proud", "rancid", "absurd", "entertaining", "annoyed", "fussy", "precise", "subtle", "gilded", "slow", "delinquent", "nervous", "hopeful", "rich", "adequate", "shrill", "plump", "freezing", "nasty", "endless", "lavish", "worried", "courteous", "bulky", "fair", "diminutive", "groggy", "miserable", "horrid", "crooked", "monstrous", "superb", "contrary", "lazy", "fidgety", "menacing", "swift", "stale", "quarrelsome", "quiet", "askew", "tough", "simple", "sweet", "hardworking", "frosty", "whispering", "famished", "crispy", "caring", "capable", "tiny", "immense", "startled", "lovely", "highpitched", "tasteless", "decrepit", "tense", "lousy", "straight", "excited", "ugly", "stunning", "parched", "wild", "ripe", "lonely", "optimistic", "obnoxious", "cavernous", "different", "harsh", "creaky", "grand", "difficult", "temporary", "eccentric", "muffled", "alert", "delicate", "timid", "infamous", "enchanting", "anxious", "humble", "edgy", "severe", "repulsive", "desolate", "sleepy", "slimy", "irritable", "vigilant", "generous", "rapid", "oldfashioned", "hilly", "easy", "righteous", "joyful", "surprised", "starving", "big", "early", "compassionate", "moody", "perpetual", "dishonest", "serious", "foolish", "soft", "old", "scared", "mighty", "trendy", "curious", "hissing", "savage", "dense", "steaming", "broad", "slick", "creative", "icy", "adorable", "slight", "terrified", "intense", "noisy", "cautious", "sizzling", "blithe", "fluttering", "faint", "delighted", "smelly", "lively", "frightened", "gauzy", "long", "strict", "bored", "calm", "melodic", "spicy", "relaxed", "triangular", "dull", "wise", "dangerous", "smooth", "cruel", "creeping", "dawdling", "intimidating", "exhausted", "deep", "tasty", "obtuse", "graceful", "tranquil", "raspy", "selfish", "sullen", "malicious", "ecstatic", "wrinkly", "opulent", "polite", "fetid", "husky", "prudent", "skinny", "tricky", "impatient", "loyal", "fresh"];
    animals = ["fawn", "peacock", "civet", "seastar", "pigeon", "bull", "bumblebee", "crocodile", "elephant", "baboon", "porcupine", "wolverine", "sparrow", "manatee", "possum", "swallow", "wildcat", "bandicoot", "labradoodle", "dragonfly", "tarsier", "chameleon", "boykin", "puffin", "bison", "llama", "kitten", "stinkbug", "macaw", "parrot", "prawn", "panther", "dogfish", "fennec", "frigatebird", "turkey", "cockatoo", "neanderthal", "crow", "gopher", "reindeer", "anaconda", "panda", "ant", "puppy", "moose", "binturong", "wildebeest", "lovebird", "ferret", "persian", "dalmatian", "bird", "umbrellabird", "kingfisher", "kangaroo", "stallion", "ostrich", "owl", "affenpinscher", "caiman", "octopus", "meerkat", "buck", "donkey", "quetzal", "chamois", "sponge", "hamster", "orangutan", "uakari", "doberman", "dormouse", "ocelot", "sparrow", "spitz", "stoat", "dragonfly", "cougar", "alligator", "walrus", "frog", "tiger", "armadillo", "chinchilla", "crab", "squid", "calf", "shrew", "dolphin", "dingo", "turtle", "chimpanzee", "armadillo", "rabbit", "basking", "coyote", "chinook", "osprey", "fly", "tiffany", "dodo", "worm", "cat", "warthog", "peccary", "shark", "pony", "monkey", "swan", "whippet", "beagle", "cougar", "anteater", "quail", "liger", "cheetah", "woodpecker", "egret", "eagle", "moose", "warthog", "snail", "budgie", "molly", "magpie", "rhinoceros", "elephant", "kudu", "wombat", "goat", "lamb", "tropicbird", "human", "hog", "tang", "lemur", "ox", "dog", "lizard", "echidna", "wallaby", "hawk", "dove", "jellyfish", "sloth", "macaque", "starfish", "guppy", "deer", "impala", "porpoise", "gazelle", "bichon", "seal", "wolf", "mole", "narwhal", "hedgehog", "sheep", "horse", "bluetick", "colt", "wildebeest", "piranha", "basenji", "mallard", "bear", "bird", "badger", "hammerhead", "kangaroo", "mule", "weasel", "dogfish", "dachsbracke", "oyster", "bat", "python", "coati", "platypus", "salamander", "cat", "caterpillar", "giraffe", "snake", "kid", "falcon", "robin", "tern", "dingo", "bolognese", "drake", "goose", "rat", "iguana", "quail", "mouse", "roebuck", "fish", "poodle", "frog", "wolverine", "chinchilla", "bobcat", "carolina", "shepherd", "snail", "mandrill", "leopard", "echidna", "rabbit", "bison", "barracuda", "foal", "ass", "eagle", "octopus", "avocet", "siamese", "dodo", "yorkie", "cockroach", "wallaroo", "tiger", "woodlouse", "fossa", "buffalo", "zorse", "albatross", "indri", "seahorse", "lemur", "louse", "ostrich", "millipede", "joey", "pinscher", "dachshund", "pelican", "chihuahua", "dogo", "wasp", "siberian", "yak", "stingray", "foxhound", "sheep", "stork", "horse", "monkey", "waterbuck", "dunker", "cuscus", "ibis", "giraffe", "aardvark", "hummingbird", "otter", "pike", "pika", "stickbug", "pelican", "dugong", "bongo", "lemming", "shrimp", "piglet", "gemsbok", "tuatara", "rottweiler", "ewe", "coati", "cichlid", "akita", "gharial", "duck", "steer", "setter", "pufferfish", "donkey", "mink", "macaw", "wolfhound", "ram", "ant", "rat", "marten", "crab", "koala", "starfish", "partridge", "chipmunk", "ibex", "maltese", "clumber", "butterfly", "flamingo", "opossum", "parrot", "mastiff", "okapi", "salmon", "tapir", "adelie", "lynx", "basilisk", "oyster", "chipmunk", "locust", "dog", "cottontop", "hyena", "oriole", "cobra", "pug", "monitor", "mandrill", "antelope", "chinstrap", "zebra", "chicken", "mule", "seal", "goat", "gull", "caterpillar", "tamarin", "wrasse", "woodchuck", "otter", "penguin", "porcupine", "bear", "ferret", "dusky", "nightingale", "bat", "jaguar", "humboldt", "ermine", "saola", "emu", "lobster", "weasel", "nightingale", "hound", "bombay", "platypus", "uguisu", "scorpion", "fox", "jerboa", "zebu", "lion", "zonkey", "ragdoll", "caracal", "bee", "kiwi", "puma", "jackal", "malamute", "mayfly", "baboon", "terrier", "jellyfish", "vicuna", "penguin", "muskrat", "zebra", "burmese", "orangutan", "himalayan", "newt", "cow", "fish", "puffin", "chin", "anteater", "beaver", "canary", "hamster", "sloth", "collie", "heron", "gopher", "magpie", "flounder", "opossum", "pademelon", "capybara", "boar", "turkey", "muskox", "bulldog", "pronghorn", "reindeer", "llama", "pygmy", "kinkajou", "cuttlefish", "cub", "bloodhound", "squirrel", "gander", "moorhen", "emu", "javanese", "birman", "harrier", "tortoise", "antelope", "gnu", "kingfisher", "wasp", "olm", "havanese", "canaan", "lizard", "ocelot", "mist", "hare", "discus", "cony", "orca", "rooster", "peacock", "akbash", "somali", "beaver", "mouse", "eland", "squirrel", "serval", "chimpanzee", "snowshoe", "toucan", "catfish", "lynx", "coyote", "bunny", "retriever", "cow", "balinese", "vulture", "coral", "leopard", "raccoon", "okapi", "kakapo", "whale", "bonobo", "moray", "cormorant", "bracke", "camel", "markhor", "rockhopper", "neapolitan", "woodpecker", "hippopotamus", "puma", "camel", "alligator", "heron", "axolotl", "argentino", "human", "mongoose", "drever", "quokka", "elk", "wombat", "civet", "panther", "gar", "lionfish", "snake", "crane", "newt", "raven", "tortoise", "chicadee", "pig", "manatee", "centipede", "numbat", "falcon", "angelfish", "chamois", "rhinoceros", "shark", "flamingo", "pheasant", "ladybird", "grasshopper", "greyhound", "lemming", "pig", "marmoset", "eel", "yorkiepoo", "mosquito", "quoll", "chick", "guanaco", "walrus", "badger", "ainu", "squid", "pekingese", "gerbil", "duck", "rattlesnake", "tapir", "lobster", "catfish", "mustang", "wallaby", "mongrel", "butterfly", "booby", "fox", "rattlesnake", "cockroach", "tadpole", "lark", "ape", "mare", "tetra", "dhole", "cesky", "raccoon", "newfoundland", "marmoset", "stag", "bullfrog", "crocodile", "lion", "barb", "wolf", "beetle", "pointer", "meerkat", "owl", "reptile", "fousek", "gibbon", "budgerigar", "swan", "hartebeest", "cassowary", "oryx", "alpaca", "gerbil", "chameleon", "vulture", "barracuda", "insect", "hen", "hare", "polecat", "fly", "yak", "seahorse", "spider", "eel", "burro", "crane", "bandicoot", "hedgehog", "dromedary", "goose", "budgerigar", "dolphin", "dormouse", "duckbill", "springbok", "mongoose", "bobcat", "gecko", "hornet", "iguana", "koala", "marmot", "skink", "deer", "filly", "barnacle", "appenzeller", "doe", "gecko", "mole", "mau", "termite", "salamander", "parakeet", "finch", "hippopotamus", "hummingbird", "cheetah", "albatross", "jaguar", "toad", "hyena", "gorilla", "skunk", "impala", "jackal", "skunk", "grouse", "moth", "caribou", "dugong", "gorilla", "chicken", "buffalo"];

    var randAnimal = 'none';
    var randAdjective = 'something';



    while (randAnimal[0] != randAdjective[0]) {
    randAnimal = animals[Math.floor(Math.random() * animals.length)];

    randAdjective = adjectives[Math.floor(Math.random() * adjectives.length)];

    }
        return toTitleCase(randAdjective) + toTitleCase(randAnimal);
}

function makeid()
{
    return makeAnimalDoublet();
    var text = "";
    var possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

    for( var i=0; i < 5; i++ )
        text += possible.charAt(Math.floor(Math.random() * possible.length));

    return text;
}

function getXMLHttp(){
    var xmlhttp;
    if (window.XMLHttpRequest) {// code for IE7+, Firefox, Chrome, Opera, Safari
        xmlhttp=new XMLHttpRequest();
        }
    else {// code for IE6, IE5
        xmlhttp=new ActiveXObject("Microsoft.XMLHTTP");
    }
    return xmlhttp;
}

function Calculate(){
    try {
        var server = document.getElementById("ML").value + "/calculate" + "?group=" + user_info.group;
        var xmlhttp = getXMLHttp();
        xmlhttp.onreadystatechange=function() {
            if (xmlhttp.readyState==4){
                //alert(xmlhttp.responseText);
            }
        }
        xmlhttp.open("GET", server);
        xmlhttp.send();
    }
    catch(err){
        //alert(err);
    }
}

function SendWifiData(network_data){
    try {
        var server = document.getElementById("ML").value;
        if (document.getElementById("learn").checked) {
            server += "/learn"
        }
        else if (document.getElementById("track").checked) {
            server += "/track"
        }
        var data = {
            "group": groupfind,
            "username": userfind,
            "password": "none",
            "location": locationfind,
            "time": Date.now(),
            "wifi-fingerprint": network_data
        }
        var xmlhttp = getXMLHttp();
        xmlhttp.onreadystatechange=function() {
            message = xmlhttp.responseText;
            if (xmlhttp.readyState==4 && xmlhttp.status == 200){
                response_json = JSON.parse(xmlhttp.responseText);
                currentLocation = response_json['position'][0][0];
                document.getElementById("currentLocation").innerHTML = 'Current location: ' + currentLocation;
            }
        }
        xmlhttp.open("POST", server);
        xmlhttp.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
        xmlhttp.send(JSON.stringify(data));
    }
    catch(err){
        //alert(err);
        TrackWifiOff();
    }
}

function successCallback(networks){
    var network_data = []
    for (var i=0; i < networks.length; i++){
        network_data.push({"mac": networks[i].BSSID,"rssi": networks[i].level})
    }
    SendWifiData(network_data);
}

function errorCallback(response){
    var output = '';
    for (var property in response) {
      output += property + ': ' + response[property]+'; ';
    }
    //alert(output);
    //TrackWifiOff();
    message = output;
}

function TrackWifiOff(){
    if (id != null) {
        document.getElementById('wifiOn').setAttribute('style','display:none;');
        document.getElementById('wifiOff').setAttribute('style','display:block;');
        navigator.wifi.clearWatch(id);
        //alert("trackig off: " + id);
        message = "trackig off: " + id;
        id = null;
        cordova.plugins.backgroundMode.configure({
                text:'FIND is running.'
        });
    }
}

function TrackWifiOn(){
    TrackWifiOff();
    if (id == null) {
        // var frequency = parseInt(document.getElementById("seconds").value) * 1000;
        document.getElementById('wifiOn').setAttribute('style','display:block;');
        document.getElementById('wifiOff').setAttribute('style','display:none;');
        userfind =  document.getElementById("userfind").value.toLowerCase();
        groupfind =  document.getElementById("groupfind").value.toLowerCase();
        locationfind =  document.getElementById("locationfind").value.toLowerCase();
        seconds =  document.getElementById("seconds").value.toLowerCase();
        window.localStorage.setItem('userfind',userfind);
        window.localStorage.setItem('groupfind',groupfind);
        window.localStorage.setItem('locationfind',locationfind);
        window.localStorage.setItem('seconds',seconds);
        window.localStorage.setItem('ML',document.getElementById("ML").value);
        var frequency = parseFloat(document.getElementById("seconds").value) * 1000;
        id = navigator.wifi.watchAccessPoints(successCallback, errorCallback, {"frequency":frequency});
        //alert("tracking on: " + id);
        message = "tracking on: " + id;
        console.log(message)

        if (document.getElementById("learn").checked) {
            cordova.plugins.backgroundMode.configure({
                    text:'Learning ' + locationfind
            });
        }
        else if (document.getElementById("track").checked) {
            cordova.plugins.backgroundMode.configure({
                    text:'Tracking'
            });
        }

    }
}

function openDashboard(){
    groupfind =  document.getElementById("groupfind").value.toLowerCase();;
    var server = document.getElementById("ML").value + "/login?group=" + groupfind;
    server = server.replace("//","/").replace('http:/','http://');
    //alert(server);
    navigator.app.loadUrl(server, { openExternal:true });
}
/*
function DisplayTime(){
    document.getElementById("time").innerHTML = document.getElementById("seconds").value;
}
*/

function showConfirm() {
        navigator.app.exitApp();
    }



var app = {
    // Application Constructor
    initialize: function() {
        this.bindEvents();
    },
    // Bind Event Listeners
    //
    // Bind any events that are required on startup. Common events are:
    // 'load', 'deviceready', 'offline', and 'online'.
    bindEvents: function() {
        document.addEventListener('deviceready', this.onDeviceReady, false);
    },
    // deviceready Event Handler
    //
    // The scope of 'this' is the event. In order to call the 'receivedEvent'
    // function, we must explicitly call 'app.receivedEvent(...);'
    onDeviceReady: function() {
        app.receivedEvent('deviceready');
        app.addListeners();
    },
    // Update DOM on a Received Event
    receivedEvent: function(id) {
        var parentElement = document.getElementById(id);
        var listeningElement = parentElement.querySelector('.listening');
        var receivedElement = parentElement.querySelector('.received');
window.localStorage.getItem('locationfind')
        listeningElement.setAttribute('style', 'display:none;');
        receivedElement.setAttribute('style', 'display:none;');


        console.log('Received Event: ' + id);
    },
    // ADD LISTENERS FOR BUTTONS
    addListeners: function(){
        document.getElementById("on").addEventListener("touchstart", TrackWifiOn);
        document.getElementById("off").addEventListener("touchstart", TrackWifiOff);
        document.getElementById("exit").addEventListener("touchstart", showConfirm);
        document.getElementById("dashboard").addEventListener("touchstart", openDashboard);
    }


};

document.getElementById('wifiOn').setAttribute('style','display:none;');

app.initialize();

document.addEventListener('pause',function() {
    // EXIT IF NOT DOING FINGERPRINTING
    if (id == null) {
    showConfirm();
    }
});


document.addEventListener('deviceready', function () {
    // Android customization
    cordova.plugins.backgroundMode.setDefaults({ text:'FIND is running.'});
    // Enable background mode
    cordova.plugins.backgroundMode.enable();

    // Called when background mode has been activated
    cordova.plugins.backgroundMode.onactivate = function () {

    }
}, false);


test = window.localStorage.getItem('userfind');
if (test == null || test.length < 1) {
    document.getElementById("userfind").value  = 'username';
} else {
    document.getElementById("userfind").value  = test;
}
test = window.localStorage.getItem('groupfind');
if (test == null || test.length < 1) {
    document.getElementById("groupfind").value  = makeid();
} else {
    document.getElementById("groupfind").value  = test;
}
test = window.localStorage.getItem('locationfind');
if (test == null || test.length < 1) {
    document.getElementById("locationfind").value  = 'your location';
} else {
    document.getElementById("locationfind").value  = test;
}
test = window.localStorage.getItem('ML');
if (test == null || test.length < 1) {
    document.getElementById("ML").value  = 'http://finddemo.duckdns.org';
} else {
    document.getElementById("ML").value  = test;
}
test = window.localStorage.getItem('seconds');
if (test == null || test.length < 1) {
    document.getElementById("seconds").value  = '2.0';
} else {
    document.getElementById("seconds").value  = test;
}
