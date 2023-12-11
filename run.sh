#!/bin/bash

image_name=$1
report_dir="reports/$image_name"
if [ ! -d "$report_dir" ]; then
    mkdir -p "$report_dir"
    echo "Created $report_dir."
fi
./trivy image --cache-dir=./cache/trivy --offline-scan --format template --template "@contrib/html.tpl" -o "$report_dir"/report.html "$image_name"