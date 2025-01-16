export async function onRequestGet(ctx) {
    try {
        // Parse the original URL and extract the pathname
        const originalUrl = new URL(ctx.request.url);
        const pathname = originalUrl.pathname;

        const rewrittenUrl = `https://cdn.akshatsharma.xyz${pathname}`;

        // Fetch the file from the MEDIA storage
        const file = await ctx.env.MEDIA.get(pathname);
        if (!file) {
            return new Response("File not found.", { status: 404 });
        }

        return new Response(file.body, {
            headers: { "Content-Type": file.httpMetadata?.contentType || "application/octet-stream" },
        });
    } catch (err) {
        return new Response(`Error: ${err.message}`, { status: 500 });
    }
}
