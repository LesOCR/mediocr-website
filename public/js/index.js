(function() {
  var $logo_highlight = $('.js-logo_highlight');


  var timeout = 0;

  $logo_highlight.hide();

  // This function is ran periodically to update the timeout and fade off the logo if needed
  setInterval(function() {
    if (timeout == 0)
      $logo_highlight.fadeOut(800);
    else if (timeout < 0)
      timeout = 0;
    else
      timeout -= 50;
  }, 50);

  // Primary event that shows the logo and sets the timeout for its disappearance
  $(document).on('mousemove', function(event) {
    if (timeout <= 0)
  	  $logo_highlight.fadeIn(500);
    timeout = 1500;
  });
})();