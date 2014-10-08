(function() {

  $(document).click(function()
  {
      animateLogo();
  });


  function animateLogo()
  {
      $animateLogo = jQuery('<div/>', {
          class: 'logo-image-animate',
          style: 'opacity: 0'
      });
      $animateLogo.appendTo('.logo-container').animate({
         opacity: 0.6,
         height: "68",
         width: "68",
         top: "50",
         left: "50"
     }, 200, function()
     {
         $animateLogo.animate({
            opacity: 0,
            height: "64",
            width: "40",
            top: "50%",
            left: "50%"
        }, 300, function()
        {
            $animateLogo.remove();
        });
     });
  }
})();
