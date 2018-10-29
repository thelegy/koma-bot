#!/bin/bash
set -ev

sass \
  --update \
  --no-source-map \
  --style "${SASS_STYLE}" \
  static/css/:static/css

postcss \
  --use autoprefixer \
  --replace \
  static/css/{,**/}*.css
