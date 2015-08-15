function newGame() {
  var templ = [
    '<h2>Make your choice</h2>',
    '<div class="box">',
      '<div>',
        '<input type="button" class="choice" name="choice" value="rock" onclick="sendChoice(\'rock\',\'{{=it.GameId}}\');">',
        '<input type="button" class="choice" name="choice" value="paper" onclick="sendChoice(\'paper\',\'{{=it.GameId}}\');">',
        '<input type="button" class="choice" name="choice" value="scissors" onclick="sendChoice(\'scissors\',\'{{=it.GameId}}\');">',
        '<input type="button" class="choice" name="choice" value="lizard" onclick="sendChoice(\'lizard\',\'{{=it.GameId}}\');">',
        '<input type="button" class="choice" name="choice" value="spock" onclick="sendChoice(\'spock\',\'{{=it.GameId}}\');">',
      '</div>',
    '</div>',
    "<div id='results'>",
    "</div>",
  ].join("\n");
  var tmp = $.post("/rpsls/new/")
  .done(function(data) {
    var response = JSON.parse(data);
    console.log(response);
    var tmp = doT.template(templ)(response);
    console.log(tmp);
    $('#content').empty().append(tmp);
  });
}

function sendChoice(selection, gameId) {
  var templ = [
      "<div class='result'>",
        "<h2><span id='result'>{{=it.LastGame}}</span></h2>",
      "</div>",
      "<h2>Game Stats</h2>",
      "<table class='table'>",
            "<tr>",
              "<th>Game ID</th>",
              "<td id='comment'>{{=it.GameId}}</td>",
            "</tr>",
            "<tr>",
              "<th>Games Played</th>",
              "<td class='comment'><span id='games'>{{=it.Games}}</span></td>",
            "</tr>",
            "<tr>",
              "<th>Player Score</th>",
              "<td class='comment'><span id='winsPlayer'>{{=it.PlayerWins}}</span></td>",
              "<td class='stat'><div class='bar'>",
              "<div class='p' style='display: in-line; width: {{=it.PlayerPercent}}%;'><span id='scorePlayer'>{{=it.PlayerPercent}}%</span></div>",
              "<div class='c' style='display: in-line; width: {{=it.ComputerPercent}}%;'><span id='scoreComputer'>{{=it.ComputerPercent}}%</span></div>",
              "<div class='t' style='display: in-line; width: {{=it.TiesPercent}}%;'><span id='scoreTies'>{{=it.TiesPercent}}%</span></div>",
              "</div></td>",
            "</tr>",
            "<tr>",
              "<th>Computer Score</th>",
              "<td class='comment'><span id='winsComputer'>{{=it.ComputerWins}}</span></td>",
            "</tr>",
            "<tr>",
              "<th>Ties</th>",
              "<td class='comment'><span id='winsTies'>{{=it.Ties}}</span></td>",
            "</tr>",
          "</table>",
      ].join("\n");
	var jqxhr = $.post( "/rpsls/result/"+gameId, {"choice": selection})
  .done(function(data) {
    var response = JSON.parse(data);
    var tmp = doT.template(templ)(response);
    $('#results').empty().append(tmp);
    $('#div:hidden:first').fadeIn("slow");
  });
}