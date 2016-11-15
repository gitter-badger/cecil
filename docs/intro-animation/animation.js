
var s = Snap(1200, 600);
scene = Snap.load("scene.svg", postLoadScene)

var numClicks = 0

function postLoadScene(loadedScene) {

    s.append(loadedScene)

    topBanner = Snap.select("#topbanner")

    
    programmer1 = Snap.select("#programmer1")
    programmer1lease = Snap.select("#programmer1lease")    
    

    programmer2 = Snap.select("#programmer2")
    programmer2lease = Snap.select("#programmer2lease")        
    
    
    robot = Snap.select("#robot")
    

    servers1 = Snap.select("#servers1")
    servers1lease = Snap.select("#servers1lease")    
    

    servers2 = Snap.select("#servers2")
    servers2lease = Snap.select("#servers2lease")    
    

    temperature = Snap.select("#temperature")
    temperature.attr({fill: "white"})

    aws = Snap.select("#aws")


    moneytemperature = Snap.select("#moneytemperature")


    hide_everything()

    next()

    // call next() every 5 seconds
    // window.setInterval(next, 2000)
    
}


function hide_everything() {
    topBanner.attr({opacity: 0.0})
    
    programmer1.attr({opacity: 0.0})
    programmer1lease.attr({opacity: 0.0})
    
    programmer2.attr({opacity: 0.0})
    programmer2lease.attr({opacity: 0.0})    

    robot.attr({opacity: 0.0})

    servers1.attr({opacity: 0.0})
    servers1lease.attr({opacity: 0.0})    

    servers2.attr({opacity: 0.0})
    servers2lease.attr({opacity: 0.0})    

    aws.attr({opacity: 0.0})		      

    moneytemperature.attr({opacity: 0.0})

    
}

function updateBanner(bannerText, fontSize) {
    topBanner.attr({opacity: 0.0})    
    topBanner.attr({ text: bannerText, style: "font-size: 3px", opacity: 0 });
    topBanner.animate({ opacity: 1 }, 500)
}

function show_the_solution() {
    updateBanner("Introducing ...", 16)
}


function show_aws() {
    aws.attr({opacity: 1.0})    
    updateBanner("You have a cloud", 16)

}

function hide_programmer_1() { 
    updateBanner("And forget about them", 16)                   
    programmer1.animate({opacity: 0.0}, 1000)
}

function hide_programmer_2() {
    updateBanner("And vacation in Costa Rica", 16)                       
    programmer2.animate({opacity: 0.0}, 1000)
}

function temperature_red() {
    updateBanner("And you waste money", 16)
    moneytemperature.animate({opacity: 1.0}, 1000)    
    temperature.attr({fill: "red"})    
}

function show_robot() {
    robot.animate({opacity: 1.0}, 2000)
}


var next_functions = [
    function() {
	show_aws()	
    },
    function() {
	updateBanner("And developers", 16)
	programmer1.animate({opacity: 1.0}, 1000)
	programmer2.animate({opacity: 1.0}, 1000)		
    },
    function() {
	updateBanner("Devs create instances", 16)            
	servers1.animate({opacity: 1.0}, 1000)
	servers2.animate({opacity: 1.0}, 1000)	
    },
    hide_programmer_1,
    hide_programmer_2,    
    temperature_red,
    function() {
	hide_everything()
	show_the_solution()
    },
    function() {
	updateBanner("Cecil - a mopster for your cloud", 16)
	show_robot()	
    },    
    function() {
	updateBanner("Devs create instances, take 2", 16)
	aws.attr({opacity: 1.0})
	robot.attr({opacity: 0.0})	
	programmer1.animate({opacity: 1.0}, 1000)
	programmer2.animate({opacity: 1.0}, 1000)		
	servers1.animate({opacity: 1.0}, 1000)
	servers2.animate({opacity: 1.0}, 1000)	
    },
    function() {
	show_robot()		
	updateBanner("Cecil creates leases w/ devs", 16)
	programmer1lease.animate({opacity: 1.0}, 1000)
	programmer2lease.animate({opacity: 1.0}, 1000)
	servers1lease.animate({opacity: 1.0}, 1000)
	servers2lease.animate({opacity: 1.0}, 1000)	
    },        
    function() {
	updateBanner("When devs vanish", 16)
	programmer1.animate({opacity: 0.0}, 1000)
	programmer2.animate({opacity: 0.0}, 1000)
	programmer1lease.animate({opacity: 0.0}, 1000)
	programmer2lease.animate({opacity: 0.0}, 1000)	
    },
    function() {
	updateBanner("Their leases expire", 16)
	servers1lease.animate({opacity: 0.0}, 1000)
	servers2lease.animate({opacity: 0.0}, 1000)
    },    
    function() {
	updateBanner("Cecil mops up instances", 16)
	servers1.animate({opacity: 0.0}, 1000)
	servers2.animate({opacity: 0.0}, 1000)
    },
    function() {
	updateBanner("And you save money", 16)
	moneytemperature.animate({opacity: 1.0}, 1000)    	
	temperature.attr({fill: "green"})
    },
    function() {
	updateBanner("And your cloud is squeaky clean", 16)
    },

]

function next() {
    f = next_functions[0]
    console.log(f)
    f()
    next_functions.shift()
}


