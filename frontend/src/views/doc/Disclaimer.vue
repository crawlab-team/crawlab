<template>
  <div class="app-container disclaimer">
    <el-card>
      <div v-html="text"></div>
    </el-card>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'
import showdown from 'showdown'

export default {
  name: 'Disclaimer',
  computed: {
    ...mapState('lang', [
      'lang'
    ]),
    text () {
      if (!this.converter) return ''
      if (this.lang === 'zh') {
        return this.converter.makeHtml(this.textZh)
      } else {
        return this.converter.makeHtml(this.textEn)
      }
    }
  },
  data () {
    const converter = new showdown.Converter()
    console.log(converter)
    return {
      converter,
      textEn: `
# Disclaimer

This Disclaimer and privacy protection statement (hereinafter referred to as "disclaimer statement" or "this statement") is applicable to the series of software (hereinafter referred to as "crawlab") developed by crawlab development group (hereinafter referred to as "development group") after you read this statement, if you do not agree with any terms in this statement or have doubts about this statement, please stop using our software immediately. If you have started or are using crawlab, you have read and agree to all terms of this statement.

1. General: by installing crawlab and using the services and functions provided by crawlab, you have agreed to establish this agreement with the development team. The developer group may at any time change the terms at its sole discretion. The amended "terms" shall take effect automatically as soon as they are published on the GitHub disclaimer page.
2. This product is a distributed crawler management platform based on golang, supporting python, nodejs, go, Java, PHP and other programming languages as well as a variety of crawler frameworks.
3. The development team of crawlab shall not be responsible for any accident, negligence, contract damage, defamation, copyright or intellectual property infringement caused by the use of crawlab and any loss caused by it (including computer virus infection caused by downloading crawlab on the unofficial site), and shall not bear any legal responsibility.
4. The user shall bear the risk of using crawlab by himself, we do not make any form of guarantee, and we will not bear any legal responsibility for the user's failure to upgrade and update normally due to any technical reasons such as network condition and communication line.
5. When users use crawlab to grab the target website, they need to comply with the laws and regulations related to crawlers, such as the network security law. Do not collect personal information of citizens without authorization, cause the target website to be paralyzed by DDoS, or fail to comply with the robots.txt protocol and other illegal means of the target website.
6. Crawlab respects and protects the personal privacy of all users and will not steal any information from users' computers.
7. Copyright of the system: the crawleb development team owns the intellectual property rights, copyrights, copyrights and use rights for all developed or jointly developed products, which are protected by applicable intellectual property rights, copyrights, trademarks, service trademarks, patents or other laws.
8. Communication: any company or individual who publishes or disseminates our software on the Internet is allowed, but the crawlab development team shall not be responsible for any legal and criminal events that may be caused by the company or individual disseminating the software.
      `,
      textZh: `
# 免责声明

本免责及隐私保护声明(下简称“免责声明”或“本声明”)适用于 Crawlab 开发组 (以下简称“开发组”)研发的系列软件(以下简称"Crawlab") 在您阅读本声明后若不同意此声明中的任何条款，或对本声明存在质疑，请立刻停止使用我们的软件。若您已经开始或正在使用 Crawlab，则表示您已阅读并同意本声明的所有条款之约定。

1. 总则：您通过安装 Crawlab 并使用 Crawlab 提供的服务与功能即表示您已经同意与开发组立本协议。开发组可随时执行全权决定更改“条款”。经修订的“条款”一经在 Github 免责声明页面上公布后，立即自动生效。
2. 本产品是基于Golang的分布式爬虫管理平台，支持Python、NodeJS、Go、Java、PHP等多种编程语言以及多种爬虫框架。
3. 一切因使用 Crawlab 而引致之任何意外、疏忽、合约毁坏、诽谤、版权或知识产权侵犯及其所造成的损失(包括在非官方站点下载 Crawlab 而感染电脑病毒)，Crawlab 开发组概不负责，亦不承担任何法律责任。
4. 用户对使用 Crawlab 自行承担风险，我们不做任何形式的保证， 因网络状况、通讯线路等任何技术原因而导致用户不能正常升级更新，我们也不承担任何法律责任。
5. 用户使用 Crawlab 对目标网站进行抓取时需遵从《网络安全法》等与爬虫相关的法律法规，切勿擅自采集公民个人信息、用 DDoS 等方式造成目标网站瘫痪、不遵从目标网站的 robots.txt 协议等非法手段。
6. Crawlab 尊重并保护所有用户的个人隐私权，不会窃取任何用户计算机中的信息。
7. 系统的版权：Crawlab 开发组对所有开发的或合作开发的产品拥有知识产权，著作权，版权和使用权，这些产品受到适用的知识产权、版权、商标、服务商标、专利或其他法律的保护。
8. 传播:任何公司或个人在网络上发布，传播我们软件的行为都是允许的，但因公司或个人传播软件可能造成的任何法律和刑事事件 Crawlab 开发组不负任何责任。
      `
    }
  }
}
</script>

<style scoped>
  .disclaimer >>> ol li {
    margin: 10px 0;
  }
</style>
