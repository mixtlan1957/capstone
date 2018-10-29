
var cytoscape = require('cytoscape');
var express = require('express');
var app = express();
var mongoose = require("mongoose");
var dbUrl = process.env.MONGODB_URI;
mongoose.connect(dbUrl, { useNewUrlParser: true});

app.use(express.static("public"));


//bodyParser!
var bodyParser = require('body-parser');
app.use(bodyParser.urlencoded({ extended: true }));
app.use(bodyParser.json());

//set view enngined
app.set('view engine', 'ejs');


const PORT = process.env.PORT || 5000;



//db stuff (work in progress)
var nodeSchema = new mongoose.Schema({
    name: String,
    children: [String],
    sqli: Boolean,
    xss: Boolean
});

var WebNode = mongoose.model("WebNode", nodeSchema);
//webNode.create;

/*
var testSite1 = new WebNode({
  name: "super.com",
  children:["bilbo.com", "reddit.com"],
  sqli: false,
  xss: false
})
*/

/*
testSite1.save(function(err, WebNode) {
  if(err) {
    console.log('Something went wrong');
  } else {
    console.log("successfully added to db");
  }
});
*/

/*
WebNode.create({
  name: "bilbo.com",
  children:[],
  sqli: false,
  xss: false 
  
}, function(err, webNode) {
  if(err) {
    console.log(err);
  } else {
    console.log("Added: " + webNode);
  }
});
*/


WebNode.find({}, function(err, webNodes) {
  if(err) {
    console.log("Error retrieving data!");
    console.log(err);
  } else {
    console.log("All the nodes:");
    console.log(webNodes);
  }
});


//create a json object
//citataion: https://stackoverflow.com/questions/36856232/write-add-data-in-json-file-using-node-js
var fs = require("fs");
var obj = [];

//this function was for testing purposes only and is to become obsolete
//written to pull data graph data from database instead of request body
function createDataForm() {
  //get what is in the database
  WebNode.find({}, function(err, webNodes) {
    if(err) {
      console.log("Error retrieving data!");
      console.log(err);
      
    //if no error retrieving from database...  
    } else {
      console.log("All the nodes:");
      
      //add nodes to JSON table
      for (var i = 0; i < webNodes.length; i++) {
        
        //only root node should initially need to be pushed
        if (i == 0) {
          obj.push({data: {id:  webNodes[i].name}});
        }
        //console.log(webNodes[i].children.length);
        for (var j = 0; j < webNodes[i].children.length; j++) {
          
          var parentNode = webNodes[i].name;
          var childName = webNodes[i].children[j];
          var concatName = parentNode + childName;
          
          //add the child node
          obj.push({data: {id: childName} } );
          
          //add the edge
          obj.push({data: {id: concatName, source: parentNode, target: childName } });
          
        }
      }
      //convert to JSON
      var json = JSON.stringify(obj);
      
      //use fs to write to file to disk
      var fs = require('fs');
      //ATTENTION: the name of the file MUST be named "data.json" for the import to work
      fs.writeFile('public/data.json', json, function(err){
        if (err) {
          throw err;
        } else {
          console.log('completed writing json file');
        }
      });
    }
  });
}

//todo: validate request input
function validateRequest(input) {

}

//This function takes a JSON request (body) and creates a JSON object
//that matches the formatting required by the 'data.json' file for
//graph generation using cytoscape.js
function createDataFromJSON(jsonInput) {
  console.log(jsonInput);

  //validate request: todo

  //create JSON output object
  var obj = [];
  //add nodes to JSON table
  for(var i = 0; i < jsonInput.length; i++) {
    //root node case
    if (i == 0 && jsonInput[i].xss == true && jsonInput[i].sqli == true) {
      obj.push({data: {id:  jsonInput[i].name, "xss": true, 
     "sqli": true, "root": true}});
    }
    else if (i == 0 && jsonInput[i].xss == true && jsonInput[i].sqli == false) {
      obj.push({data: {id:  jsonInput[i].name, "xss": true, 
      "sqli": false, "root": true}});
    }
    else if (i == 0 && jsonInput[i].xss == false && jsonInput[i].sqli == true) {
      obj.push({data: {id:  jsonInput[i].name, "xss": false, 
      "sqli": true, "root": true}});
    }
    else if (i == 0) {
      obj.push({data: {id:  jsonInput[i].name, "xss": false, 
      "sqli": false, "root": true}});
    }

    //all other nodes
    if (i > 0 && jsonInput[i].xss == true && jsonInput[i].sqli == true) {
     obj.push({data: {id:  jsonInput[i].name, "xss": true, 
     "sqli": true, "root": false}});
    }
    else if(i > 0 && jsonInput[i].xss == true && jsonInput[i].sqli == false) {
      obj.push({data: {id:  jsonInput[i].name, "xss": true, 
      "sqli": false, "root": false}});
    }
    else if(i > 0 && jsonInput[i].xss == false && jsonInput[i].sqli == true) {
      obj.push({data: {id:  jsonInput[i].name, "xss": false, 
      "sqli": true, "root": false}});
    }
    else if (i > 0) {
      obj.push({data: {id:  jsonInput[i].name, "xss": false, 
      "sqli": false, "root": false}});
    }

    //add edges
    for (var j = 0; j < jsonInput[i].children.length; j++) {

      var parentNode = jsonInput[i].name;
      var childName = jsonInput[i].children[j];
      var concatName = parentNode + childName;  //define the edge ie connection
      
      //add the edge
      obj.push({data: {id: concatName, source: parentNode, target: childName } });
    }
  }
  //convert to JSON
  var jsonOutput = JSON.stringify(obj);

  return jsonOutput;
}

//create data before page is loaded (we could also use a another callback)
//createDataForm();

//routes
app.get('/', function(req, res){
  WebNode.find({}, function(err, all_webNodes){
    if(err) {
      console.log(err);
    } else {
      //createDataForm();
      res.render('pages/index'); 
    }
  });
});

app.post('/graphs', function(req, res){

  var graphData = createDataFromJSON(req.body);
  //console.log(graphData);

  fs.writeFile('public/data.json', graphData, function(err){
    if(err) {
      throw err;
    } else {
      console.log('post request successfully written to file');
      res.status(201);
      res.redirect('/');
    }
  });
});




//404
app.use(function(req,res){
  res.status(404);
  res.render('pages/404');
});


//505
app.use(function(err, req, res, next){
  console.error(err.stack);
  res.type('plain/text');
  res.status(500);
  res.render('pages/500');
});



app.listen(PORT, process.env.IP, function() {
  //var addr = app.address();
  //console.log("App started at ", addr.address + ":" + addr.port);
  console.log("App started");
});

