<?php
$file = $argv[1];

list($width, $height) = getimagesize($file);

$image = new Imagick();
$image->readImage($file);

$image->setFirstIterator();

$frames = [];

$image = $image->coalesceImages();
do {
    $delay = $image->getImageDelay();
    $image->setImageFormat('jpg');
    $frames[] = [
        'data' => base64_encode($image->getImageBlob()),
        'delay' => $delay
    ];
} while ($image->nextImage());

echo sprintf('<svg xmlns="http://www.w3.org/2000/svg" xmlns:A="http://www.w3.org/1999/xlink" width="%d" height="%d">', $width, $height) . PHP_EOL;

foreach ($frames as $k => $frame) {
    echo sprintf('  <image id="%06d" height="100%%" A:href="data:image/jpeg;base64,' . $frame['data'] . '"/>', $k) . PHP_EOL;
}

foreach ($frames as $k => $frame) {
    if ($k == 0) {
        $begin = sprintf('A%06d.end; 0s', count($frames) - 1);
    } else {
        $begin = sprintf('A%06d.end', $k - 1);
    }
    echo sprintf('  <set A:href="#%06d" id="A%06d" attributeName="width" to="100%%" dur="%dms" begin="%s"/>', $k, $k, $frame['delay'] * 10, $begin) . PHP_EOL;
}

echo '</svg>' . PHP_EOL;
