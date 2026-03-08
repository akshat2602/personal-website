export async function onRequestGet(ctx) {
  const url = new URL(ctx.request.url);
  const key = url.pathname.slice(1);

  const object = await ctx.env.MEDIA.get(key);

  if (!object) {
    return new Response("Not found", { status: 404 });
  }

  const headers = new Headers();
  object.writeHttpMetadata(headers);
  headers.set("Cache-Control", "public, max-age=31536000, immutable");

  return new Response(object.body, { headers });
}
