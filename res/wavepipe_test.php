<?php

require_once __DIR__ . "/wavepipe.php";

// Check for host argument
if (empty($argv[1])) {
	printf("No host provided");
	exit(1);
}
$host = $argv[1];

// Attempt a login request using test credentials
$login = json_decode(@file_get_contents("http://" . $host . "/api/v0/login?u=test&p=test"), true);
if (empty($login)) {
	printf("Failed to decode login JSON");
	exit(1);
}

// Store necessary login information
$publicKey = $login["session"]["publicKey"];
$secretKey = $login["session"]["secretKey"];

// Iterate and test all JSON APIs
$apiCalls = array(
	"/api/v0/albums",
	"/api/v0/albums/1",
	"/api/v0/artists",
	"/api/v0/artists/1",
	"/api/v0/folders",
	"/api/v0/folders/1",
	"/api/v0/search/song",
	"/api/v0/songs",
	"/api/v0/songs/1",
	"/api/v0/logout",
);

foreach ($apiCalls as $a) {
	// Create a nonce
	$nonce = generateNonce();

	// Create the necessary API signature
	$signature = apiSignature($publicKey, $nonce, "GET", $a, $secretKey);

	// Generate URL
	$url = sprintf("http://%s%s?s=%s:%s:%s", $host, $a, $publicKey, $nonce, $signature);

	printf("%s:\n", $a);

	// Perform API call
	$contents = @file_get_contents($url);
	if (empty($contents)) {
		printf("Failed to retrieve data stream");
		exit(1);
	}

	printf("%s\n", $contents);
}
