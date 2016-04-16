var scanningInterval;
var toggle = true;
var currentLocation = "none";
var isMovile;
var press;
var groupname;
var username;
var servername;


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

function openDashboard() {
 var servername = window.localStorage.getItem("server").toLowerCase();
    if (servername.slice(-1) != '/') {
    	servername += "/";
    }
    if (servername.indexOf("http") < 0) {
    	servername = "http://" + servername;
    }
    var groupname = window.localStorage.getItem("group").toLowerCase();
	navigator.app.loadUrl(servername + "login?group=" + groupname, { openExternal:true });
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

function ipIntToString(ip) {
	var str = "";
	for(var i=0; i<4; i++) {
		str += ip % 256;
		ip = Math.floor(ip / 256);
		if(i < 3) str += '.';
	}
	return str;
}


function learn() {
	console.log('button#scanwifi');
	if (toggle==true) {
		navigator.notification.prompt(
		    'Please enter location name',  // message
		    scanAndSend,                  // callback to invoke
		    'Learn location',            // title
		    ['Ok','Exit'],             // buttonLabels
		    'office'                 // defaultText
		);	
		toggle = false;
	}
}


function sendFingerprint() {
	if(window.plugins && window.plugins.WifiAdmin) {
	    // Enable background mode
	    if (cordova.plugins.backgroundMode.isEnabled() == false) {
		    cordova.plugins.backgroundMode.enable();	    	
	    }

		var wf = window.plugins.WifiAdmin;
		wf.getWifiInfo(function(data){
		console.log( JSON.stringify(data) );
		
		var wifiConnected = data['activity'];
		var wifiList = data['available'];
		
		var html = "";
		if(wifiConnected != null) {
			html += "Connected to:<br/>" +
				"SSID: " + wifiConnected['SSID'] + "<br/>" +
				"BSSID: " + wifiConnected['BSSID'] + "<br/>" + 
				"Mac Addr: " + wifiConnected['MacAddress'] + "<br/>" + 
				"IP: " + ipIntToString( wifiConnected['IpAddress'] ) + "<br/>" +
				"Speed: " + wifiConnected['LinkSpeed'] + " Mbps<br/>"; 
		} else {
			html += "Not connected.<br/>";
		}
		
		html += "<br/>Available Wifi:<br/>";
		network_data = []
		while(wifiList.length >0) {
			var item = wifiList.shift();
			html += item['BSSID'] + '(' + item['level'] + 'dB, feq:' + (item['frequency']/1000.0).toFixed(2) + 'GHz)';
			if(item['BSSID'] === wifiConnected['BSSID']) html += '(*)';
			html += '<br/>';
			network_data.push({"mac": item['BSSID'],"rssi": item['level']})
		}
		
		var data = {
            "group": window.localStorage.getItem("group").toLowerCase(),
            "username": window.localStorage.getItem("username").toLowerCase(),
            "password": "none",
            "location": currentLocation,
            "time": Date.now(),
            "wifi-fingerprint": network_data
        }
        var route = "learn";
        if (currentLocation == "tracking") {
        	route = "track";
			$('div#sending').html("Tracking current location.");
        } else {
			$('div#sending').html("Learning fingerprint for " + currentLocation );
        }
        var servername = window.localStorage.getItem("server").toLowerCase();
        if (servername.slice(-1) != '/') {
        	servername += "/";
        }
        if (servername.indexOf("http") < 0) {
        	servername = "http://" + servername;
        }
		$.ajax({
		   type: "POST",
		   url: servername + route,
		   dataType: "json",
		   data: JSON.stringify(data),
		   success: function(data) {
		   	var d = new Date();
			var n = d.toString();
		     $('div#result').html( n + "<br>" + data["message"] );
		   },
		   error: function(e) {
		     $('div#result').html('Error: ' + e.message);
		   }
		});
					
		}, function(){});
	}
}

function scanAndSend(results) {
	if (results == null) {
		results = {input1:"tracking",buttonIndex:1};
	}
	currentLocation = results.input1.toLowerCase();
	if (results.buttonIndex == 1) {
		clearInterval(scanningInterval);
		 var servername = window.localStorage.getItem("server").toLowerCase();
            if (servername.slice(-1) != '/') {
            	servername += "/";
            }
            if (servername.indexOf("http") < 0) {
            	servername = "http://" + servername;
            }
		$('div#scanning').html("Sending fingerprint to " + servername);
		sendFingerprint();
		scanningInterval = setInterval(sendFingerprint,3000);
	}
	toggle = true;
}

function stopScanning() {
	$('div#scanning').html("Not scanning");
	$('div#result').html("");
	$('div#sending').html("");

	clearInterval(scanningInterval);
    if (cordova.plugins.backgroundMode.isEnabled() == true) {
	    cordova.plugins.backgroundMode.disable();	    	
    }
}


function setData(datatype,defaultname) {
	navigator.notification.prompt(
		    'Please enter a ' + datatype + ' name',  // message
		    function(results) {
				if (results.buttonIndex == 1) {
					window.localStorage.setItem(datatype,results.input1);
				}
				$('h2#user').html("Group: " + window.localStorage.getItem("group") + "<br>User: " + window.localStorage.getItem("username"));
			},                  // callback to invoke
		    'Set ' + datatype,            // title
		    ['Ok','Exit'],             // buttonLabels
		    defaultname                 // defaultText
	);	
}



function initUIEvents() {
	var isMobile = ( /(android|ipad|iphone|ipod)/i.test(navigator.userAgent) );
	var press = isMobile ? 'touchstart' : 'mousedown';
	



	$('button#openwifi').on(press, function(){
		if(window.plugins && window.plugins.WifiAdmin) {
			var wf = window.plugins.WifiAdmin;
			wf.enableWifi(true);
		}
	});

	$('button#closewifi').on(press, function(){
		if(window.plugins && window.plugins.WifiAdmin) {
			var wf = window.plugins.WifiAdmin;
			wf.enableWifi(false);
		}
	});


}

function main() {
	initUIEvents();

	// local storage
	test = window.localStorage.getItem('username');
	if (test == null || test.length < 1) {
	    setData('username','user1');
	} 

	test = window.localStorage.getItem('group');
	if (test == null || test.length < 1) {
	    window.localStorage.setItem('group',makeid())
	} 

	test = window.localStorage.getItem('server');
	if (test == null || test.length < 1) {
	    servername  = 'https://ml.internalpositioning.com';
	    window.localStorage.setItem('server',servername)
	} 

	$('h2#user').html("Group: " + window.localStorage.getItem("group") + "<br>User: " + window.localStorage.getItem("username"));

    // Android customization
    cordova.plugins.backgroundMode.setDefaults({ text:'FIND is running.'});
    cordova.plugins.backgroundMode.disable();
    // Called when background mode has been activated
    cordova.plugins.backgroundMode.onactivate = function () {

    }

}

function onDeviceReady() {
    console.log(navigator.notification);
}
