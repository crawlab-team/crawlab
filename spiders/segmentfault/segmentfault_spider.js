const puppeteer = require('puppeteer');

(async () => {
  const browser = await (puppeteer.launch({
    timeout: 15000
  }));

  const url = 'https://segmentfault.com/newest';

  const page = await browser.newPage();

  await page.goto(url);
  await page.waitFor(2000);

  await page.screenshot({path: 'screenshot.png'});

  const titles = await page.evaluate(sel => {
    let results = [];
    document.querySelectorAll('.news-list .news-item .news__item-title').forEach(el => {
      results.push({
        title: el.innerText
      })
    });
    return results;
  });

  console.log(titles);

  browser.close();
})();