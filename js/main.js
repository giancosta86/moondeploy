$(document).ready(function() {
  var downloadButton = $("#downloadProgram")

  var lowerCasePlatform = navigator.platform.toLowerCase()

  if (lowerCasePlatform.indexOf("win") > -1) {
    var os = "windows"
  } else if (lowerCasePlatform.indexOf("linux") > -1) {
    var os = "linux"
  } else if (lowerCasePlatform.indexOf("mac") > -1) {
    var os = "mac"
  }

  println("OS = " + os)

  var osFileRegex = new RegExp("-" + os + "\.(?:exe|zip)$" )

  $.get('https://api.github.com/repos/giancosta86/moondeploy/releases/latest', function (data) {
    var osFileUrl

    data.assets.forEach(function(asset) {
      var matchesOS = osFileRegex.test(asset.name)

      if (matchesOS) {
        osFileUrl = asset.browser_download_url
      }

      println(asset.name + " --> " + matchesOS)
    })

    if (osFileUrl) {
      downloadButton.attr("href", osFileUrl)
    }
  });
})


function println(line) {
  if (console && console.log) {
    console.log(line)
  }
}
