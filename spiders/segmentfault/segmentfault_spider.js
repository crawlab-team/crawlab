const puppeteer = require('puppeteer');
const MongoClient = require('mongodb').MongoClient;

(async () => {
  // browser
  const browser = await (puppeteer.launch({
    timeout: 15000
  }));

  // define start url
  const url = 'https://segmentfault.com/newest';

  // start a new page
  const page = await browser.newPage();

  // navigate to url
  await page.goto(url);
  await page.waitFor(2000);

  // take a screenshot
  await page.screenshot({path: 'screenshot.png'});

  // scrape data
  const results = await page.evaluate(() => {
    let results = [];
    document.querySelectorAll('.news-list .news-item .news__item-title').forEach(el => {
      results.push({
        title: el.innerText
      })
    });
    return results;
  });

  // open database connection
  const client = await MongoClient.connect('mongodb://localhost/crawlab_test');
  let db = await client.db('test');
  const colName = process.env.CRAWLAB_COLLECTION;
  const taskId = process.env.CRAWLAB_TASK_ID;
  const col = db.collection(colName);

  // save to database
  await results.forEach(d => {
    d.task_id = taskId;
    col.save(d);
  });

  // close database connection
  db.close();

  console.log(results);

  // shutdown browser
  browser.close();
})();