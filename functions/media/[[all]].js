export async function onRequestGet(ctx) {
    // Parse the original URL and extract the pathname
    const originalUrl = new URL(ctx.request.url);
    const pathname = originalUrl.pathname;

    // Rewrite the URL with the new base URL
    const rewrittenUrl = new URL(pathname, "https://cdn.akshatsharma.xyz");

    // Fetch the file from the rewritten URL
    const file = await ctx.env.MEDIA.get(rewrittenUrl.pathname.replace("/media/", ""));
    if (!file) return new Response(null, { status: 404 });

    return new Response(file.body, {
        headers: { "Content-Type": file.httpMetadata.contentType },
    });
}
