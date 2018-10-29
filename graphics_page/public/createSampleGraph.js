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
            'text-valign': 'center',
            'font-size': 4
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
            'content': 'data(id)',
            'background-color': 'red',
            'shape': 'diamond'
          }
        },

        {
          selector: '.hoverOverSQLI',
          style: {
            'content': 'data(id)',
            'background-color': 'red',
            'shape': 'triangle'
          }
        },

        {
          selector: '.hoverOver',
          style: {
            'content': 'data(id)',
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
cy.on('tap', 'node', function() {
  var url = this.data('id'),
    substring = "///";

  //strip out tripple '/' if necessary
  var foundIdx = url.indexOf(substring);
  if (foundIdx !== -1) {
    var temp = url.substring(foundIdx + 3, url.length);
    url = temp;
  }



  
  try {
    //ciatation:
    //https://stackoverflow.com/questions/29684740/javascript-window-open-without-http
    if (!url.match(/^http?:\/\//i) || !url.match(/^https?:\/\//i)) {
      url = 'http://' + url;
    }
    window.open(url);
  } catch(e) {
    window.location.href = this.data(url);
  }
});


});


  
