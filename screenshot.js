const { chromium } = require('playwright');
const path = require('path');

async function screenshots() {
    const browser = await chromium.launch({
        executablePath: `${process.env.HOME}/.cache/ms-playwright/chromium-1208/chrome-linux64/chrome`,
        headless: true,
    });
    const ctx = await browser.newContext({ viewport: { width: 1280, height: 800 } });
    const page = await ctx.newPage();

    const pages = [
        ['/', 'home'],
        ['/posts/', 'archives'],
        ['/posts/i-am-not-writing-enough/', 'single-post'],
        ['/posts/experience_mongo_gen_ai_hackathon/', 'sidenotes-post'],
        ['/tags/', 'tags'],
        ['/search/', 'search'],
        ['/now/', 'now'],
    ];

    for (const [url, name] of pages) {
        await page.goto('http://localhost:8080' + url);
        await page.waitForTimeout(300);
        await page.screenshot({ path: `/tmp/ssg-${name}.png`, fullPage: true });
        console.log(`Screenshot saved: /tmp/ssg-${name}.png`);
    }

    // Mobile view of sidenote post
    await ctx.close();
    const mobileCtx = await browser.newContext({ viewport: { width: 390, height: 844 } });
    const mobilePage = await mobileCtx.newPage();
    await mobilePage.goto('http://localhost:8080/posts/experience_mongo_gen_ai_hackathon/');
    await mobilePage.waitForTimeout(300);
    await mobilePage.screenshot({ path: '/tmp/ssg-sidenotes-mobile.png', fullPage: false });
    console.log('Screenshot saved: /tmp/ssg-sidenotes-mobile.png');

    await browser.close();
    console.log('All screenshots done.');
}

screenshots().catch(e => { console.error(e); process.exit(1); });
