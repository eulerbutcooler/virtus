import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import ThemeToggle from '../../components/ui/ThemeToggle.jsx'
import './LandingPage.css'

const stats = [
  { value: '$0', label: 'Pooled Today' },
  { value: '0', label: 'Requests Fulfilled' },
  { value: '100%', label: 'Publicly Auditable' },
]

const steps = [
  {
    num: '01',
    title: 'Submit a Request',
    body: 'Describe what you need and why. An admin reviews it before it enters the queue — no anonymous gatekeeping.',
  },
  {
    num: '02',
    title: 'Community Funds It',
    body: 'Contributors add to a shared pool. When your request is covered, it moves from queue to procurement automatically.',
  },
  {
    num: '03',
    title: 'We Procure & Deliver',
    body: 'The admin sources the item, ships it, and records proof of delivery. Every step is visible on the Transparency page.',
  },
  {
    num: '04',
    title: 'You Report Impact',
    body: 'After delivery, you log how it helped. That record stays on-chain with your request — closed-loop accountability.',
  },
]

const fadeUp = {
  hidden: { opacity: 0, y: 24 },
  visible: (i = 0) => ({
    opacity: 1,
    y: 0,
    transition: { duration: 0.55, ease: [0.22, 1, 0.36, 1], delay: i * 0.08 },
  }),
}

const stagger = {
  visible: { transition: { staggerChildren: 0.1 } },
}

export default function LandingPage() {
  return (
    <div className="landing">

      <nav className="landing__nav">
        <Link to="/" className="landing__wordmark">Virtus</Link>
        <div className="landing__nav-links">
          <Link to="/transparency" className="landing__nav-link">Transparency</Link>
          <Link to="/login" className="landing__nav-link">Sign In</Link>
          <ThemeToggle />
          <Link to="/register" className="landing__nav-cta">Get Started</Link>
        </div>
      </nav>


      <motion.section
        className="landing__hero"
        initial="hidden"
        animate="visible"
        variants={stagger}
      >
        <motion.div className="landing__hero-eyebrow" variants={fadeUp}>
          Community Fulfillment Platform
        </motion.div>
        <motion.h1 className="landing__hero-title" variants={fadeUp} custom={1}>
          Real needs.<br />
          Real community.<br />
          <span className="landing__hero-accent">Full transparency.</span>
        </motion.h1>
        <motion.p className="landing__hero-sub" variants={fadeUp} custom={2}>
          Virtus is a mutual-aid platform where every contribution, request, and delivery is
          publicly auditable. No black boxes. No middlemen skimming fees.
        </motion.p>
        <motion.div className="landing__hero-actions" variants={fadeUp} custom={3}>
          <Link to="/register" className="landing__btn-primary">Join the Community</Link>
          <Link to="/transparency" className="landing__btn-ghost">View Live Data →</Link>
        </motion.div>


        <motion.div className="landing__stats" variants={stagger}>
          {stats.map((s, i) => (
            <motion.div key={s.label} className="landing__stat" variants={fadeUp} custom={i + 4}>
              <span className="landing__stat-value">{s.value}</span>
              <span className="landing__stat-label">{s.label}</span>
            </motion.div>
          ))}
        </motion.div>
      </motion.section>


      <motion.section
        className="landing__how"
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true, margin: '-80px' }}
        variants={stagger}
      >
        <motion.div className="landing__section-header" variants={fadeUp}>
          <span className="landing__section-eyebrow">How It Works</span>
          <h2 className="landing__section-title">From request to delivery — nothing hidden</h2>
        </motion.div>

        <div className="landing__steps">
          {steps.map((step, i) => (
            <motion.div
              key={step.num}
              className="landing__step"
              variants={fadeUp}
              custom={i}
            >
              <div className="landing__step-num">{step.num}</div>
              <div className="landing__step-content">
                <h3 className="landing__step-title">{step.title}</h3>
                <p className="landing__step-body">{step.body}</p>
              </div>
              {i < steps.length - 1 && <div className="landing__step-connector" />}
            </motion.div>
          ))}
        </div>
      </motion.section>


      <motion.section
        className="landing__values"
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true, margin: '-60px' }}
        variants={stagger}
      >
        {[
          { icon: '⬡', title: 'Open Pool', body: 'One shared pool. Every contribution goes directly toward the next request in queue — proportional, fair, verifiable.' },
          { icon: '◎', title: 'Zero Opacity', body: 'Every transaction, approval, and delivery is logged and publicly visible at /transparency — no login required.', offset: true },
          { icon: '◈', title: 'Closed Loop', body: 'Impact records close the cycle. Contributors see exactly what their money built, with photos and outcome notes.' },
        ].map((v, i) => (
          <motion.div
            key={v.title}
            className={['landing__value-card', v.offset ? 'landing__value-card--offset' : ''].filter(Boolean).join(' ')}
            variants={fadeUp}
            custom={i}
          >
            <div className="landing__value-icon">{v.icon}</div>
            <h3>{v.title}</h3>
            <p>{v.body}</p>
          </motion.div>
        ))}
      </motion.section>


      <motion.section
        className="landing__cta"
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true, margin: '-60px' }}
        variants={stagger}
      >
        <motion.h2 className="landing__cta-title" variants={fadeUp}>
          Ready to contribute<br />or submit a request?
        </motion.h2>
        <motion.p className="landing__cta-sub" variants={fadeUp} custom={1}>
          It takes two minutes to join. Your first contribution starts the cycle.
        </motion.p>
        <motion.div variants={fadeUp} custom={2}>
          <Link to="/register" className="landing__btn-primary landing__btn-primary--lg">
            Create an Account
          </Link>
        </motion.div>
      </motion.section>


      <footer className="landing__footer">
        <Link to="/" className="landing__wordmark">Virtus</Link>
        <div className="landing__footer-links">
          <Link to="/transparency">Transparency</Link>
          <Link to="/login">Sign In</Link>
          <Link to="/register">Register</Link>
        </div>
        <span className="landing__footer-copy">Built for people, by people.</span>
      </footer>
    </div>
  )
}
