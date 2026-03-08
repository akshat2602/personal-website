/**
 * Playwright tests for the custom SSG-generated site.
 * Run after starting a local server: python3 -m http.server 8080 --directory public
 */
const { chromium } = require('playwright');

const BASE = 'http://localhost:8080';
let passed = 0;
let failed = 0;
const errors = [];

function assert(condition, message) {
    if (condition) {
        console.log(`  ✓ ${message}`);
        passed++;
    } else {
        console.error(`  ✗ FAIL: ${message}`);
        failed++;
        errors.push(message);
    }
}

async function runTests() {
    const browser = await chromium.launch({
        executablePath: `${process.env.HOME}/.cache/ms-playwright/chromium-1208/chrome-linux64/chrome`,
        headless: true,
    });
    const ctx = await browser.newContext();
    const page = await ctx.newPage();

    // ── 1. Home page ────────────────────────────────────────────────────────
    console.log('\n[1] Home page');
    await page.goto(BASE + '/');
    assert(await page.title() === 'Akshat Sharma', 'Title is "Akshat Sharma"');
    assert(await page.locator('.header .nav .logo a').innerText() === 'Akshat Sharma', 'Logo text correct');
    assert(await page.locator('#menu').isVisible(), 'Navigation menu visible');
    assert(await page.locator('.home-info').isVisible(), 'Home info section visible');
    assert(await page.locator('.home-info h1').innerText() !== '', 'Home info has heading');
    assert(await page.locator('.social-icons').isVisible(), 'Social icons visible');
    // Posts should be listed
    const postEntries = await page.locator('.post-entry').count();
    assert(postEntries > 0, `Post entries visible on home (found ${postEntries})`);
    assert(postEntries <= 5, `Pagination: ≤5 posts per page (got ${postEntries})`);
    // CSS loaded — check background is dark Dracula
    const bgColor = await page.evaluate(() =>
        window.getComputedStyle(document.body).backgroundColor
    );
    assert(bgColor === 'rgb(40, 42, 54)', `Dark Dracula background applied (got: ${bgColor})`);

    // ── 2. Navigation links ─────────────────────────────────────────────────
    console.log('\n[2] Navigation links');
    const navLinks = await page.locator('#menu a').all();
    assert(navLinks.length >= 4, `Nav has ≥4 links (got ${navLinks.length})`);
    const navTexts = await page.locator('#menu a').allInnerTexts();
    assert(navTexts.some(t => t.includes('posts')), 'Nav has posts link');
    assert(navTexts.some(t => t.includes('tags')), 'Nav has tags link');
    assert(navTexts.some(t => t.includes('search')), 'Nav has search link');
    assert(navTexts.some(t => t.includes('now')), 'Nav has now link');

    // ── 3. Pagination ───────────────────────────────────────────────────────
    console.log('\n[3] Pagination');
    const nextLink = page.locator('.pagination a.next');
    assert(await nextLink.isVisible(), 'Pagination "Older" link visible');
    await nextLink.click();
    assert(page.url().includes('/page/2/'), 'Page 2 URL correct');
    const page2Posts = await page.locator('.post-entry').count();
    assert(page2Posts > 0, `Page 2 has posts (${page2Posts})`);
    const prevLink = page.locator('.pagination a.prev');
    assert(await prevLink.isVisible(), 'Page 2 has "Newer" link');

    // ── 4. Single post page ─────────────────────────────────────────────────
    console.log('\n[4] Single post page');
    await page.goto(BASE + '/posts/i-am-not-writing-enough/');
    assert(await page.title() === 'I am not writing enough | Akshat Sharma', 'Post title correct');
    assert(await page.locator('.post-title').isVisible(), 'Post title visible');
    assert(await page.locator('.post-content').isVisible(), 'Post content visible');
    assert(await page.locator('.post-meta').isVisible(), 'Post meta visible');
    const metaText = await page.locator('.post-meta').innerText();
    assert(metaText.includes('min read'), 'Reading time shown');
    assert(metaText.includes('words'), 'Word count shown');
    // Breadcrumbs
    assert(await page.locator('.breadcrumbs').isVisible(), 'Breadcrumbs visible');
    const breadcrumbText = await page.locator('.breadcrumbs').innerText();
    assert(breadcrumbText.includes('Home') && breadcrumbText.includes('Posts'), 'Breadcrumbs content correct');
    // Tags
    assert(metaText.includes('writing'), 'Tags visible in meta');
    // Bluesky comments section (has bsky frontmatter)
    assert(await page.locator('#comments-section').isVisible(), 'Bluesky comments section present');
    const bskyUri = await page.locator('#comments-section').getAttribute('data-bsky-uri');
    assert(bskyUri !== null && bskyUri.includes('bsky.app'), `Bluesky URI set correctly: ${bskyUri}`);
    // Post nav links
    assert(await page.locator('.paginav').isVisible(), 'Post navigation visible');
    // Share buttons
    assert(await page.locator('.share-buttons').isVisible(), 'Share buttons visible');

    // ── 5. Post with sidenotes ──────────────────────────────────────────────
    console.log('\n[5] Post with sidenotes');
    await page.goto(BASE + '/posts/experience_mongo_gen_ai_hackathon/');
    assert(await page.title() !== '', 'Sidenote post loads');
    const sidenotes = await page.locator('.sidenote').count();
    assert(sidenotes >= 2, `Sidenotes rendered (found ${sidenotes})`);
    const sidenoteLabels = await page.locator('.sidenote-label').count();
    assert(sidenoteLabels >= 2, `Sidenote labels rendered (found ${sidenoteLabels})`);
    const sidenoteContent = await page.locator('.sidenote-content').count();
    assert(sidenoteContent >= 2, `Sidenote content spans rendered (found ${sidenoteContent})`);
    // On desktop the sidenote-right should be positioned absolute
    const sidenoteRightPos = await page.evaluate(() => {
        const el = document.querySelector('.sidenote-content.sidenote-right');
        if (!el) return null;
        return window.getComputedStyle(el).position;
    });
    assert(sidenoteRightPos === 'absolute', `Sidenote-right is absolute positioned (got: ${sidenoteRightPos})`);

    // ── 6. Table of Contents ────────────────────────────────────────────────
    console.log('\n[6] Table of Contents');
    await page.goto(BASE + '/posts/i-am-not-writing-enough/');
    const toc = page.locator('.toc');
    assert(await toc.isVisible(), 'TOC visible on post with headings');
    const tocLinks = await toc.locator('a').count();
    assert(tocLinks >= 2, `TOC has links (found ${tocLinks})`);

    // ── 7. Archives page ────────────────────────────────────────────────────
    console.log('\n[7] Archives page');
    await page.goto(BASE + '/posts/');
    assert(await page.title() === 'Posts | Akshat Sharma', 'Archives page title correct');
    assert(await page.locator('.archive-year').count() > 0, 'Year groups present');
    const archiveMonths = await page.locator('.archive-month').count();
    assert(archiveMonths > 0, `Archive months present (${archiveMonths})`);
    const archiveEntries = await page.locator('.archive-entry').count();
    assert(archiveEntries > 0, `Archive entries present (${archiveEntries})`);
    // Click an archive entry
    const firstEntry = page.locator('.archive-entry a.entry-link').first();
    const entryHref = await firstEntry.getAttribute('href');
    assert(entryHref !== null && entryHref.startsWith('/posts/'), `Archive entry href correct: ${entryHref}`);

    // ── 8. Tags pages ───────────────────────────────────────────────────────
    console.log('\n[8] Tags pages');
    await page.goto(BASE + '/tags/');
    assert(await page.title() === 'Tags | Akshat Sharma', 'Tags page title');
    const tagItems = await page.locator('.terms-tags li').count();
    assert(tagItems > 0, `Tags listed (${tagItems})`);
    // Click a tag
    await page.locator('.terms-tags li a').first().click();
    assert(page.url().includes('/tags/'), 'Tag page URL correct');
    const tagPosts = await page.locator('.post-entry').count();
    assert(tagPosts > 0, `Tag page shows posts (${tagPosts})`);

    // ── 9. Search page ──────────────────────────────────────────────────────
    console.log('\n[9] Search page');
    await page.goto(BASE + '/search/');
    assert(await page.title() === 'Search | Akshat Sharma', 'Search page title');
    assert(await page.locator('#searchInput').isVisible(), 'Search input visible');
    // Wait for fuse.js search index to load
    await page.waitForFunction(() => window.searchReady === true, { timeout: 5000 });
    // Type a search term
    await page.locator('#searchInput').fill('raft');
    await page.locator('#searchInput').dispatchEvent('keyup');
    await page.waitForTimeout(300);
    const results = await page.locator('#searchResults li').count();
    assert(results > 0, `Search returns results for "raft" (${results})`);

    // ── 10. Now page ────────────────────────────────────────────────────────
    console.log('\n[10] Now page');
    await page.goto(BASE + '/now/');
    assert(await page.title() === 'Now | Akshat Sharma', 'Now page title');
    assert(await page.locator('.page-header h1').isVisible(), 'Now page heading visible');

    // ── 11. RSS feed ────────────────────────────────────────────────────────
    console.log('\n[11] RSS feed');
    await page.goto(BASE + '/index.xml');
    const rssContent = await page.content();
    assert(rssContent.includes('<rss'), 'RSS feed contains <rss> element');
    assert(rssContent.includes('<title>'), 'RSS feed has title');
    assert(rssContent.includes('<item>'), 'RSS feed has items');

    // ── 12. Search JSON index ───────────────────────────────────────────────
    console.log('\n[12] Search JSON index');
    await page.goto(BASE + '/index.json');
    const jsonText = await page.evaluate(() => document.body.innerText);
    let searchIndex;
    try {
        searchIndex = JSON.parse(jsonText);
        assert(Array.isArray(searchIndex), 'Search index is a JSON array');
        assert(searchIndex.length > 0, `Search index has entries (${searchIndex.length})`);
        assert(searchIndex[0].title !== undefined, 'Search entries have title');
        assert(searchIndex[0].permalink !== undefined, 'Search entries have permalink');
    } catch(e) {
        assert(false, `Search index is valid JSON: ${e.message}`);
    }

    // ── 13. CSS / Syntax highlighting ───────────────────────────────────────
    console.log('\n[13] CSS & Syntax highlighting');
    // Find a post with code
    await page.goto(BASE + '/posts/cockroachdb/');
    const hasChroma = await page.locator('.chroma').count();
    // Not all posts have code, but CSS should load
    const cssLoaded = await page.evaluate(() =>
        getComputedStyle(document.body).getPropertyValue('--link-hover-color').trim()
    );
    assert(cssLoaded === '#6272a4', `CSS dark theme variable correct (--link-hover-color: "${cssLoaded}")`);

    // ── 14. 404 page ────────────────────────────────────────────────────────
    console.log('\n[14] 404 page');
    await page.goto(BASE + '/404.html');
    assert(await page.locator('.not-found').isVisible(), '404 page has .not-found element');

    // ── 15. OG images generated ─────────────────────────────────────────────
    console.log('\n[15] OG images');
    const ogResp = await page.goto(BASE + '/og/i-am-not-writing-enough.png');
    assert(ogResp.status() === 200, `OG image accessible for "i-am-not-writing-enough" (status: ${ogResp.status()})`);

    // ── 16. Sidenote mobile toggle ──────────────────────────────────────────
    console.log('\n[16] Sidenote mobile toggle');
    await page.setViewportSize({ width: 600, height: 800 });
    await page.goto(BASE + '/posts/experience_mongo_gen_ai_hackathon/');
    // On mobile, sidenote-right should be hidden by default
    const mobileDisplay = await page.evaluate(() => {
        const el = document.querySelector('.sidenote-content.sidenote-right');
        if (!el) return null;
        return window.getComputedStyle(el).display;
    });
    assert(mobileDisplay === 'none', `Sidenote hidden on mobile by default (got: ${mobileDisplay})`);
    // Click label to toggle
    const label = page.locator('.sidenote-label').first();
    await label.click();
    const mobileDisplayAfter = await page.evaluate(() => {
        const el = document.querySelector('.sidenote-content.sidenote-right');
        if (!el) return null;
        return window.getComputedStyle(el).display;
    });
    assert(mobileDisplayAfter !== 'none', `Sidenote shows after click on mobile (got: ${mobileDisplayAfter})`);
    // Reset viewport
    await page.setViewportSize({ width: 1280, height: 800 });

    // ── 17. Scroll-to-top button ─────────────────────────────────────────────
    console.log('\n[17] Scroll-to-top button');
    await page.goto(BASE + '/posts/i-am-not-writing-enough/');
    const topLink = page.locator('#top-link');
    assert(await topLink.count() > 0, 'Scroll-to-top link exists');

    // ── 18. External links open in new tab ──────────────────────────────────
    console.log('\n[18] External nav links');
    await page.goto(BASE + '/');
    const resumeLink = page.locator('#menu a[href*="resume"]');
    const resumeTarget = await resumeLink.getAttribute('target');
    assert(resumeTarget === '_blank', 'Resume link opens in new tab');

    await browser.close();

    // ── Summary ──────────────────────────────────────────────────────────────
    console.log('\n' + '─'.repeat(50));
    console.log(`Results: ${passed} passed, ${failed} failed`);
    if (errors.length > 0) {
        console.error('\nFailed assertions:');
        errors.forEach(e => console.error(`  - ${e}`));
    }
    process.exit(failed > 0 ? 1 : 0);
}

runTests().catch(err => {
    console.error('Test runner error:', err);
    process.exit(1);
});
