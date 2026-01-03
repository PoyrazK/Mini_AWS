const payload = JSON.parse(process.env.PAYLOAD || '{}');
console.log('Hello from CloudFunctions!');
console.log('Payload Received:', JSON.stringify(payload));
if (payload.error) {
    process.exit(1);
}
