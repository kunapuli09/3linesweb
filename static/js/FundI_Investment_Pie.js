
var ctxP = document.getElementById("pieChart").getContext('2d');
  var myPieChart = new Chart(ctxP, {
    type: 'pie',
    data: {
      labels: ["Artificial Intelligence", "Enterprise Software", "Cyber Security", "Clean Tech", "Fin Tech", "Supply Chain", "Mobile Apps", ],
      datasets: [{
        data: [22, 13, 5,5, 46, 5,6],
        backgroundColor: ["#7FDBFF", "#2ECC40", "#FFDC00", "#39CCCC", "#FF851B","#FF4136", "#0074D9"],
        hoverBackgroundColor: ["#FF5A5E", "#5AD3D1", "#FFC870", "#A8B3C5", "#616774"]
      }]
    },
    options: {
      responsive: true
    }
  });


  //line
  var ctxL = document.getElementById("lineChart").getContext('2d');
  var myLineChart = new Chart(ctxL, {
    type: 'line',
    data: {
      labels: ["Q3 2017", "Q4 2017", "Q1 2018", "Q2 2018", "Q3 2018", "Q4 2018", "Q1 2019"],
      datasets: [{
          label: "Invested Amount",
          data: [50000, 100000, 415000, 1615000, 2240000, 2340000, 2590000],
          backgroundColor: ["#1E90FF",
          ],
          borderColor: [
            'rgba(200, 99, 132, .7)',
          ],
          borderWidth: 2
        },
        {
          label: "Portfolio Value",
          data: [62500, 125000, 518750, 2018750, 2800000, 2925000, 3237500],
          backgroundColor: ["#2CEC2F",
          ],
          borderColor: [
            'rgba(0, 10, 130, .7)',
          ],
          borderWidth: 2
        }
      ]
    },
    options: {
      responsive: true
    }
  });

