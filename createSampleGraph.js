//load root node


fetch('data.json', {mode: 'no-cors'}).then(function(res) {
  return res.json()
}).then(function(data) {
  
var cy = window.cy = cytoscape({
    container: document.getElementById('cy'),
 
  
    boxSelectionEnabled: false,
    autounselectify: true,
 
    layout: {
        name: 'cose-bilkent',
        animate: false
      }, 
    
    style: [
        {
          selector: 'node',
          style: {
            'background-color': 'blue',
            'width': 6,
            'height': 6,
            'text-valign': 'top',
            'font-size': 3
            //'content': 'data(id)'
          }
        },
        {
          selector: '.sqliFound',
          style: {
            'background-color': 'blue',
            'shape': 'triangle'
          }
        },

        {
          selector: '.xssFound',
          style: {
            'background-color': 'blue',
            'shape': 'diamond'
          }
        },

        {
          selector: 'edge',
          style: {
            'width': .5,
            'line-color': '#ad1a66'
          }
        },

        {
          selector: '.hoverOverXSS',
          style: {
            'content': 'data(material)',
            'background-color': 'red',
            'shape': 'diamond'
          }
        },

        {
          selector: '.hoverOverSQLI',
          style: {
            'content': 'data(material)',
            'background-color': 'red',
            'shape': 'triangle'
          }
        },

        {
          selector: '.hoverOver',
          
          style: {
            'content': 'data(material)',
            'background-color': 'blue'
          }
        }

      ],
      
  
    
      
    //load subsequent nodes
    elements: data

      
    //change the style of the nodes if necessary

   
    
});
//changes the format based on wether or not xss or sqli flag has been marked
cy.nodes().forEach(function( ele ){
  if (ele.data("sqli") == true){
    //cy.$(ele.data("id")).classes("secondClass");
    ele.classes('sqliFound');
    ele.animate({
      style: {
        'background-color': 'red'
      }
    });
  } else if (ele.data("xss") == true) {
    ele.classes('xssFound');
    ele.animate({
      style: {
        'background-color': 'red'
      }
    });
  } else if (ele.data("root") == true) {
    ele.animate({
      style: {
        'background-color': 'pink'
      }
    });
  } else if (ele.data("keyword") == true) {
    ele.animate({
      style: {
        'background-color': 'green'
      }
    });
  }
});



//mouseover effects (wip)

cy.on('mouseover', 'node', function(evt) {
  var node = evt.target;

  if (node.data("sqli") == true) {
    node.classes('hoverOverSQLI');
  }
  else if (node.data("xss") == true) {
    node.classes('hoverOverXSS');
  }
  else {
    node.classes('hoverOver');
  }

  cy.on('mouseout', 'node', function(evt){
    var node = evt.target;
    if (node.data("sqli") == true) {
      node.classes('sqliFound');
    }
    else if (node.data("xss") == true) {
      node.classes('xssFound');
    }
    else {
      node.classes('node');
    }
  });
});

//follow link
cy.on('tap', 'node', function(evt) {

  var node = evt.target;
  // console.log('tapped' + node.id()) // gives id
  console.log(node); //object itself

  // console.log("node._private");

  // console.log(node._private);

  // console.log(node._private.edges);
  // console.log(node._private.edges.length);

  var nodeData = this.data()

  //dummy arrays for what would hold the real fuzz results
  nodeData.testInfo = [ "input with id 'user' failed sql test",
  "input with id 'name' failed sql test"];

  nodeData.XssTestInfo = [ "reflection xss detected" ];

  //make table
  var newTable = document.createElement("table");
  newTable.style.border = "1px solid blue";
  newTable.style.position = "absolute";
  
  newTable.style.top = 200 + "px"; //could make x and y
  newTable.style.left = 200 + "px"; //could make this x and y

  document.body.appendChild(newTable);

  //td-tr
  for (var i = 0; i < 2; i++) {
        var newTr = document.createElement("tr");
        var newTd = document.createElement("td");

        newTr.appendChild(newTd);
        newTable.appendChild(newTr);

        newTd.style.width = "380px"; //dimensions of the display
        newTd.style.height = "43px";
        newTd.style.background = "lightgrey";

        //banner
        if (i === 0) {
            newTd.style.color = "white";
            var newSpan = document.createElement("span");
            newTd.appendChild(newSpan);
            newSpan.textContent = "[X]";
            newSpan.style.backgroundColor = "red";
            newSpan.addEventListener("click", function() {
                var parentTable = this.parentNode.parentNode.parentNode;
                parentTable.parentNode.removeChild(parentTable);
            })

            // so that only the x is clickable
            newSpan = document.createElement("span");
            newSpan.textContent = nodeData.material + " statistics";

            newTd.appendChild(newSpan);
            newTd.style.background = "blue";

          }

          //body
          else {
             newTd.style.height = "275px";
             newTd.style.verticalAlign = "top";
             newTd.innerHTML = 
             `<p> <strong> Seed site: </strong> <a href="${nodeData.id}"> ${nodeData.id} </a> </p>
              <p> <strong> Number of edges: </strong> ${node._private.edges.length} </p>
             `

              // {% if nodeData.sqli == false %} <p style="color:green"> SQL injection vulnerabilities: ${nodeData.sqli} </p> 
              // {% else %} <p style="color: red"> SQL injection vulnerabilities: ${nodeData.sqli} </p> {% endif %}

             if ( nodeData.keyword == false ) {
                newTd.innerHTML += `<p> <strong> Keyword: </strong> <span style="color:red"> ${nodeData.keyword} </span> </p>`;
             } else { 
                newTd.innerHTML += `<p> <strong> Keyword: </strong> <span style="color:green"> ${nodeData.keyword} </span> </p>`;
             }   

             if ( nodeData.sqli == false ) {
                newTd.innerHTML += `<p> <strong> SQL injection vulnerabilities: </strong> <span style="color:green"> ${nodeData.sqli} </span> </p>`;
             } else { 
                newTd.innerHTML += `<p> <strong> SQL injection vulnerabilities: </strong> <span style="color:red"> ${nodeData.sqli} </span> </p>`;
             } 

            if ( nodeData.xss == false ) {
                newTd.innerHTML += `<p> <strong> XSS vulnerabilities: </strong> <span style="color:green"> ${nodeData.xss} </span> </p>`;
             } else { 
                newTd.innerHTML += `<p> <strong> XSS injection vulnerabilities: </strong> <span style="color:red"> ${nodeData.xss} </span> </p>`;
             } 

            newTd.innerHTML +=
            `<span id="showResults" style="background-color:#ddd;" 
                onclick="">
                Click to show results
            </span>`;

             document.getElementById("showResults").onclick = function() {

                console.log(this);
                this.innerHTML = ""; //clear
                this.innerHTML += `<h3> SQL </h3>`;
                for (var test in nodeData.testInfo) {
                    this.innerHTML += 
                    `
                    <p> ${nodeData.testInfo[test]} </p>
                    `;
                }

                this.innerHTML += `<h3> XSS </h3>`;
                for (var test in nodeData.XssTestInfo) {
                    this.innerHTML += 
                    `
                    <p> ${nodeData.XssTestInfo[test]} </p>
                    `;
                }
             }

          }

    }


  // console.log(this.data());

  // id: "http://localhost:8000/page1.html"
  // keyword: false
  // material: "Page 1 - http://localhost:8000/page1.html"
  // root: false
  // sqli: false
  // title: "Page 1"
  // xss: false

  // console.log(node)
  // alert(this.data('id'));

  // var url = this.data('id'),
  //   substring = "///";

  // //strip out tripple '/' if necessary
  // var foundIdx = url.indexOf(substring);
  // if (foundIdx !== -1) {
  //   var temp = url.substring(foundIdx + 3, url.length);
  //   url = temp;
  // }



  
  // try {
  //   //ciatation:
  //   //https://stackoverflow.com/questions/29684740/javascript-window-open-without-http
  //   if (!url.match(/^http?:\/\//i) || !url.match(/^https?:\/\//i)) {
  //     url = 'http://' + url;
  //   }
  //   window.open(url);
  // } catch(e) {
  //   window.location.href = this.data(url);
  // }
});


});


  
