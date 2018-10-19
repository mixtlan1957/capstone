
var cytoscape = require('cytoscape');
var express = require('express');
var bodyParser = require('body-parser');
var app = express();
var mongoose = require("mongoose");
mongoose.connect("mongodb://localhost:27017/testsdb", { useNewUrlParser: true});

app.use(express.static("public"));


//bodyParser!
var bodyParser = require('body-parser');
app.use(bodyParser.urlencoded({ extended: true }));
app.use(bodyParser.json());

//set view enngined
app.set('view engine', 'ejs');


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
  name: "reddit.com",
  parent: "super.com",
  children:["newegg.com", "mgtgoldfish.com"],
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

//traverse the data stored in the 
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
        
        //only parent node should initially need to be pushed
        if (i == 0) {
          obj.push({data: {id:  webNodes[i].name}});
        }
        console.log(webNodes[i].children.length);
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

//create data before page is loaded (we could also use a another callback)
createDataForm();

//routes
app.get('/', function(req, res){
  WebNode.find({}, function(err, all_webNodes){
    if(err) {
      console.log(err);
    } else {
      
      
      res.render('pages/index');
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



app.listen(process.env.PORT, process.env.IP, function() {
  //var addr = app.address();
  //console.log("App started at ", addr.address + ":" + addr.port);
  console.log("App started");
});

