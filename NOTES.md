For bundling on Mac OS, we'll probably use: https://gist.github.com/mholt/11008646c95d787c30806d3f24b2c844

For properly setting up the exe icon and other info, we'll probably use: https://github.com/josephspurrier/goversioninfo

Both of these will probably warrant a structural change such as:

  * arbitray/
    * go/
      * go files
      * icon/
        * iconwin.go, etc. (using '../../resources/icon')
    * resources/
      * .ico
      * .png
      * etc
    * Makefile? Some sort of build steps for:
      * Making Mac OS .app
        * Code signing
      * Packaging Mac OS .app into .dmg
      * Fixing .exe for icon and other info on Windows

It might also be possible to move arbitray-js into a `js` subdirectory, but we'll see.
