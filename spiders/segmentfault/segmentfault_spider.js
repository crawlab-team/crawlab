const puppeteer = require('puppeteer');
const MongoClient = require('mongodb').MongoClient;

(async () => {
  // browser
  const browser = await (puppeteer.launch({
    headless: true
  }));

  // define start url
  const url = 'https://segmentfault.com/newest';

  // start a new page
  const page = await browser.newPage();

  // navigate to url
  try {
    await page.goto(url, {waitUntil: 'domcontentloaded'});
    await page.waitFor(2000);
  } catch (e) {
    console.error(e);

    // close browser
    browser.close();

    // exit code 1 indicating an error happened
    code = 1;
    process.emit("exit ");
    process.reallyExit(code);

    return
  }

  // scroll down to fetch more data
  for (let i = 0; i < 10; i++) {
    console.log('Pressing PageDown...');
    await page.keyboard.press('PageDown', 200);
    await page.waitFor(500);
  }

  // scrape data
  const results = await page.evaluate(() => {
    let results = [];
    document.querySelectorAll('.news-list .news-item').forEach(el => {
      results.push({
        url: 'https://segmentfault.com' + el.querySelector('.news__item-info > a').getAttribute('href'),
        title: el.querySelector('.news__item-title').innerText
      })
    });
    return results;
  });

  // open database connection
  const client = await MongoClient.connect('mongodb://127.0.0.1:27017');
  let db = await client.db('crawlab_test');
  const colName = process.env.CRAWLAB_COLLECTION || 'results_segmentfault';
  const taskId = process.env.CRAWLAB_TASK_ID;
  const col = db.collection(colName);

  // save to database
  for (let i = 0; i < results.length; i++) {
    // de-duplication
    const r = await col.findOne({url: results[i]});
    if (r) continue;

    // assign taskID
    results[i].task_id = taskId;
    results[i].source = 'segmentfault';

    // insert row
    await col.insertOne(results[i]);
  }

  console.log(`results.length: ${results.length}`);

  // close database connection
  client.close();

  // shutdown browser
  browser.close();
})();