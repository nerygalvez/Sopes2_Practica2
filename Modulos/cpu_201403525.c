/**
*
*	Creando el modulo de cpu
*
*/



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
//#include <linux/cputime.h>
#include <linux/cpuidle.h>
#include <linux/tick.h>
#include <asm/apic.h>

#include <linux/smp.h>
#include <linux/timex.h>
#include <linux/string.h>
#include <linux/seq_file.h>
#include <linux/cpufreq.h>
#include <linux/delay.h>


/////////////////////////
#include <linux/list.h>
#include <linux/types.h>
#include <asm/uaccess.h> 
#include <linux/sched/signal.h>
 

#define FileProc "201403525_cpu"
#define Carne "201403525"
#define Nombre "Nery Gonzalo Galvez Gomez"
#define SO "Ubuntu 18.04.4 LTS"
#define Curso "Sistemas Operativos 2"





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
		//idle = usecs_to_cputime64(idle_time);
		idle = nsecs_to_jiffies64(idle_time);

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
		//iowait = usecs_to_cputime64(iowait_time);
		iowait = nsecs_to_jiffies64(iowait_time);

	return iowait;
}

#endif








static int proc_llenar_archivo(struct seq_file *m, void *v) {
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

        //Creo el json con los datos
        seq_printf(m, "{ \"Total\" : %lld , \"Utilizado\" : %lld , \"Libre\" : %lld , \"Promedio\" : %lld }"
        , jiffies_64_to_clock_t(sum), jiffies_64_to_clock_t(sum - idle), jiffies_64_to_clock_t(idle), jiffies_64_to_clock_t(((sum - idle) * 100 / sum)));
	
	return 0;
}


static int proc_al_abrir_archivo(struct inode *inode, struct  file *file) {
  return single_open(file, proc_llenar_archivo, NULL);
}

static struct file_operations myops =
{
        .owner = THIS_MODULE,
        .open = proc_al_abrir_archivo,
        .read = seq_read,
        .llseek = seq_lseek,
        .release = single_release,
};



/**
*	Defino que es lo que se va a hacer al cargar el modulo
*/
static int iniciar(void)
{
	proc_create(FileProc,0,NULL,&myops);
    	printk(KERN_INFO "Carne: %s\n", Carne);

        /*
         * Si no se devuelve 0 significa que initmodule ha fallado y no ha podido cargarse.
         */
        return 0;
}

/**
*	Defino que es lo que se va a hacer al terminar el modulo
*/
static void terminar(void)
{
	remove_proc_entry(FileProc,NULL);
  	printk(KERN_INFO "Curso: %s\n", Curso);
}



/*
 * Indicamos cuales son las funciones de inicio y fin
 */
module_init(iniciar);
module_exit(terminar);

/*
 * Documentacion del modulo
 */
MODULE_LICENSE("GPL");
MODULE_AUTHOR("Nery Galvez - 201403525");
MODULE_DESCRIPTION("Modulo con informacion del CPU");
MODULE_SUPPORTED_DEVICE("TODOS");
