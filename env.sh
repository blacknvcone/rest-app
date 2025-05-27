#!/usr/bin/env bash

export LIBRARY_PATH="${pkgs.tesseract}/lib:${pkgs.leptonica}/lib:${pkgs.opencv4}/lib"
export CPATH="${pkgs.tesseract}/include:${pkgs.leptonica}/include:${pkgs.opencv4}/include"
export PKG_CONFIG_PATH="${pkgs.tesseract}/lib/pkgconfig:${pkgs.leptonica}/lib/pkgconfig:${pkgs.opencv4}/lib/pkgconfig"
# GoCV (OpenCV for Go)
export CGO_CPPFLAGS="-I${pkgs.opencv4}/include/opencv4 -I${pkgs.opencv4}/include"
export CGO_CXXFLAGS="--std=c++11"
export CGO_LDFLAGS="-L${pkgs.opencv4}/lib -lopencv_core -lopencv_highgui -lopencv_imgproc -lopencv_videoio -lopencv_imgcodecs"
