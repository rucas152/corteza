<!doctype html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<!-- Bootstrap core CSS -->
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@4.6.0/dist/css/bootstrap.min.css" rel="stylesheet">
		<!-- Bootstrap icons -->
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.3.0/font/bootstrap-icons.css">
		<!-- Custom CSS -->
		<link href="{{ links.Assets }}/style.css?{{ buildtime }}" rel="stylesheet">
	</head>
	<body style="background: url({{ links.Assets }}/background.png);">
		<main class="auth">
			<div class="card">
			{{ template "inc_nav.html.tpl" . }}