/************************************************
* SE UTILIZÓ CODIGO DE REFERENCIA DE GITHUB Y STACKOVERFLOW
* https://github.com/01org/KVMGT-kernel/blob/master/fs/proc/stat.c
* AUTOR: @paulgortmaker
* https://stackoverflow.com/questions/9229333/how-to-get-overall-cpu-usage-e-g-57-on-linux?utm_medium=organic&utm_source=google_rich_qa&utm_campaign=google_rich_qa
* 
*************************************************/
#include <linux/fs.h>
#include <linux/hugetlb.h>
#include <linux/init.h>
#include <linux/kernel.h>
#include <linux/sched.h>
#include <linux/module.h>
#include <linux/seq_file.h>
#include <linux/proc_fs.h>

#include <linux/mm.h>
#include <linux/mman.h>
#include <linux/mmzone.h>

#include <linux/quicklist.h>

#include <linux/swap.h>
#include <linux/vmstat.h>
#include <linux/atomic.h>
#include <asm/page.h>
#include <asm/pgtable.h>

#include <linux/kernel.h>
#include <linux/module.h>
#include <linux/cpumask.h>
#include <linux/fs.h>
#include <linux/init.h>
#include <linux/interrupt.h>
#include <linux/kernel_stat.h>
#include <linux/proc_fs.h>
#include <linux/sched.h>
#include <linux/seq_file.h>
#include <linux/slab.h>
#include <linux/time.h>
#include <linux/irqnr.h>
#include <linux/cputime.h>
#include <linux/tick.h>
#include <asm/apic.h>

#include <linux/smp.h>
#include <linux/timex.h>
#include <linux/string.h>
#include <linux/seq_file.h>
#include <linux/cpufreq.h>
#include <linux/delay.h>

#ifdef arch_idle_time

static cputime64_t get_idle_time(int cpu)
{
	cputime64_t idle;
	idle = kcpustat_cpu(cpu).cpustat[CPUTIME_IDLE];
	if (cpu_online(cpu) && !nr_iowait_cpu(cpu))
		idle += arch_idle_time(cpu);
	return idle;
}
static cputime64_t get_iowait_time(int cpu)
{
	cputime64_t iowait;
	iowait = kcpustat_cpu(cpu).cpustat[CPUTIME_IOWAIT];
	if (cpu_online(cpu) && nr_iowait_cpu(cpu))
		iowait += arch_idle_time(cpu);
	return iowait;
}
#else

static u64 get_idle_time(int cpu)
{
	u64 idle, idle_time = -1ULL;
	if (cpu_online(cpu))
		idle_time = get_cpu_idle_time_us(cpu, NULL);
	if (idle_time == -1ULL)
		/* !NO_HZ or cpu offline so we can rely on cpustat.idle */
		idle = kcpustat_cpu(cpu).cpustat[CPUTIME_IDLE];
	else
		idle = usecs_to_cputime64(idle_time);

	return idle;
}

static u64 get_iowait_time(int cpu)
{
	u64 iowait, iowait_time = -1ULL;
	if (cpu_online(cpu))
		iowait_time = get_cpu_iowait_time_us(cpu, NULL);

	if (iowait_time == -1ULL)
		/* !NO_HZ or cpu offline so we can rely on cpustat.iowait */
		iowait = kcpustat_cpu(cpu).cpustat[CPUTIME_IOWAIT];
	else
		iowait = usecs_to_cputime64(iowait_time);

	return iowait;
}

#endif

static int meminfo_proc_show(struct seq_file *m, void *v)
{
	/**
	 * DECLARAR VARIABLES
	 * */
	int i;
	unsigned long jif;
	u64 user, nice, system, idle, iowait, irq, softirq, steal;
	u64 guest, guest_nice;
	u64 sum = 0;
	struct timespec boottime;

	user = nice = system = idle = iowait =
		irq = softirq = steal = 0;
	guest = guest_nice = 0;
	getboottime(&boottime);
	jif = boottime.tv_sec;
	//Recolectar informaciòn de cada CPU y aumentar las variables/contadores
	for_each_possible_cpu(i)
	{
		user += kcpustat_cpu(i).cpustat[CPUTIME_USER];
		nice += kcpustat_cpu(i).cpustat[CPUTIME_NICE];
		system += kcpustat_cpu(i).cpustat[CPUTIME_SYSTEM];
		idle += get_idle_time(i);
		iowait += get_iowait_time(i);
		irq += kcpustat_cpu(i).cpustat[CPUTIME_IRQ];
		softirq += kcpustat_cpu(i).cpustat[CPUTIME_SOFTIRQ];
		steal += kcpustat_cpu(i).cpustat[CPUTIME_STEAL];
		guest += kcpustat_cpu(i).cpustat[CPUTIME_GUEST];
		guest_nice += kcpustat_cpu(i).cpustat[CPUTIME_GUEST_NICE];
	}
	//El total del cpu es la suma de todos los atributos
	sum += user + nice + system + idle + iowait + irq + softirq + steal + guest + guest_nice;
	//crear el json, al ser de tipo u64, se debe utilizar cputime64_to_clock_t
	seq_printf(m, "{");
	seq_printf(m, "\"cpu\":");
	seq_put_decimal_ull(m, ' ', cputime64_to_clock_t(sum));
	seq_printf(m, ",\"used\":");
	seq_put_decimal_ull(m, ' ', cputime64_to_clock_t(sum - idle));
	seq_printf(m, ",\"free\":");
	seq_put_decimal_ull(m, ' ', cputime64_to_clock_t(idle));
	seq_printf(m, ",\"average\":");
	seq_put_decimal_ull(m, ' ', cputime64_to_clock_t(((sum - idle) * 100 / sum)));
	seq_printf(m, "}");
	return 0;
}

static void __exit final(void) //Salida de modulo
{
	printk(KERN_INFO "Cleaning Up.\n");
}

static int meminfo_proc_open(struct inode *inode, struct file *file)
{
	return single_open(file, meminfo_proc_show, NULL);
}

static const struct file_operations meminfo_proc_fops = {
	.open = meminfo_proc_open,
	.read = seq_read,
	.llseek = seq_lseek,
	.release = single_release,
};
static int __init inicio(void) //Escribe archivo en /proc
{
	proc_create("m_cpu_201443728", 0, NULL, &meminfo_proc_fops);
	return 0;
}
module_init(inicio);
module_exit(final);
MODULE_LICENSE("GPL");

