$(document).ready(function() {
  var downloadButton = $("#download-program")

  var lowerCasePlatform = navigator.platform.toLowerCase()

  if (lowerCasePlatform.indexOf("win") > -1) {
    var os = "windows"
  } else if (lowerCasePlatform.indexOf("linux") > -1) {
    var os = "linux"
  } else if (lowerCasePlatform.indexOf("mac") > -1) {
    var os = "mac"
  }

  var osFileRegex = new RegExp("-" + os + "\.(?:exe|zip)$" )

  $.get('https://api.github.com/repos/giancosta86/moondeploy/releases/latest', function (data) {
    var osFileUrl

    data.assets.forEach(function(asset) {
      if (osFileRegex.test(asset.name)) {
        osFileUrl = asset.browser_download_url
      }
    })

    if (osFileUrl) {
      downloadButton.attr("href", osFileUrl)
    }
  });
})
