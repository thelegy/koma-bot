#!/bin/bash
set -ev

sass \
  --update \
  --scss \
  --force \
  --sourcemap=none \
  --style "${SASS_STYLE}" \
  static/css

postcss \
  --use autoprefixer \
  --replace \
  static/css/**/*.css
