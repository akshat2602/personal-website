export async function onRequestGet(ctx) {
    try {
        // Parse the original URL and extract the pathname
        const originalUrl = new URL(ctx.request.url);
        const pathname = originalUrl.pathname;
        console.log(`URL path: ${pathname}`);

        const objectKey = pathname.startsWith("/") ? pathname.slice(1) : pathname;

        console.log(`Fetching object with key: ${objectKey}`);
        const object = await ctx.env.MEDIA.get(objectKey);
        if (!object) {
            return new Response("File not found in R2.", { status: 404 });
        }

        // Prepare the response with the appropriate Content-Type header
        return new Response(object.body, {
            headers: {
                "Content-Type": object.httpMetadata?.contentType || "application/octet-stream",
                "Cache-Control": "public, max-age=31536000, immutable", // Cache for a year, suitable for static assets
            },
        });
    } catch (err) {
        // Handle unexpected errors
        return new Response(`Error: ${err.message}`, { status: 500 });
    }
}
