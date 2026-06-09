<template>
  <section class="dashboard-page">
    <header class="dashboard-page-head">
      <div>
        <p class="section-kicker">Overview</p>
        <h2>项目首页</h2>
      </div>
      <p class="workspace-note">mysurl1 控制台</p>
    </header>

    <section class="intro-hero">
      <div class="intro-copy">
        <p class="eyebrow">Project Intro</p>
        <h3>一个面向登录用户的短链服务，强调清晰归属、快速创建和稳定跳转。</h3>
        <p>
          当前前端已经接入注册、登录和短链创建流程。用户登录后，通过 Bearer token 调用创建接口，短链归属与当前账号绑定。
        </p>
      </div>

      <div class="intro-grid">
        <article class="intro-card">
          <span class="signal-title">认证</span>
          <strong>Login / Register</strong>
          <p>登录和注册接口已经接入，未登录用户会被拦截到登录页。</p>
        </article>
        <article class="intro-card">
          <span class="signal-title">短链创建</span>
          <strong>POST /api/v1/links</strong>
          <p>创建请求自动携带 `Authorization: Bearer &lt;token&gt;`。</p>
        </article>
        <article class="intro-card">
          <span class="signal-title">用户隔离</span>
          <strong>User Scoped</strong>
          <p>同一长链在不同用户下可独立归属，前端只允许登录后操作。</p>
        </article>
      </div>
    </section>
  </section>
</template>
