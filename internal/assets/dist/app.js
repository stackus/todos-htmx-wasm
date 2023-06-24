htmx.onLoad(function (content) {
  var sortables = content.querySelectorAll(".sortable");
  for (var i = 0; i < sortables.length; i++) {
    var sortable = sortables[i];
    new Sortable(sortable, {
      draggable: '.draggable',
      animation: 150,
      chosenClass: 'dragClass'
    });
  }
});
