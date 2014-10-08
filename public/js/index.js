(function() {

  $(document).click(function()
  {
      animateLogo();
  });


  function animateLogo()
  {
      $animateLogo = jQuery('<div/>', {
          class: 'logo-image-animate',
      });
      $animateLogo.appendTo('.logo-container').animate({
         opacity: 0,
         height: "168",
         width: "168",
         top: "0",
         left: "0"
     }, 1000, function()
     {
         $animateLogo.remove();
     });
  }
})();
