const puppeteer = require('puppeteer');
const MongoClient = require('mongodb').MongoClient;

(async () => {
  // browser
  const browser = await (puppeteer.launch({
    headless: true
  }));

  // page
  const page = await browser.newPage();

  // open database connection
  const client = await MongoClient.connect('mongodb://192.168.99.100:27017');
  let db = await client.db('crawlab_test');
  const colName = process.env.CRAWLAB_COLLECTION || 'results';
  const col = db.collection(colName);
  const col_src = db.collection('results');

  const results = await col_src.find({content: {$exists: false}}).toArray();
  for (let i = 0; i < results.length; i++) {
    let item = results[i];

    // define article anchor
    let anchor;
    if (item.source === 'juejin') {
      anchor = '.article-content';
    } else if (item.source === 'segmentfault') {
      anchor = '.article';
    } else if (item.source === 'csdn') {
      anchor = '#content_views';
    } else {
      continue;
    }

    console.log(`anchor: ${anchor}`);

    // navigate to the article
    try {
      await page.goto(item.url, {waitUntil: 'domcontentloaded'});
      await page.waitFor(2000);
    } catch (e) {
      console.error(e);
      continue;
    }

    // scrape article content
    item.content = await page.$eval(anchor, el => el.innerHTML);

    // save to database
    await col.save(item);
    console.log(`saved item: ${JSON.stringify(item)}`)
  }

  // close mongodb
  client.close();

  // close browser
  browser.close();

})();