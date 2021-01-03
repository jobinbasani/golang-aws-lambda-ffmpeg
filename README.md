# Media Transcoder Using Container Images in AWS Lambda
This is a small example that makes use of the container image support of AWS Lambda functions.
This embeds a static ffmpeg binary in the base Go image for AWS Lambda, and transcodes the input file to WebM format.
