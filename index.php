<?php
// AddType image/svg+xml xng

header('Content-type: image/svg+xml');

echo file_get_contents('test.xng');
