using System;
using Microsoft.EntityFrameworkCore;
using Dwn.Data.Database;
using Dwn.Data.Models.Identity;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Configuration;
using Microsoft.AspNetCore.Identity;
using System.Linq;
using Dwn.Cli.Commands;
using System.Collections.Generic;

namespace Dwn.Cli
{
    class Program
    {
        private static IConfigurationRoot Configuration;
        private static ServiceCollection Services;

        static void Main(string[] args)
        {
            var provider = GetServiceProvider();

            var db = provider.GetService<ApplicationDbContext>();
            var userManager = provider.GetService<UserManager<ApplicationUser>>();
            var roleManager = provider.GetService<RoleManager<ApplicationRole>>();

            if (args.Length == 1)
            {
                Console.WriteLine("Available commands: list");
                return;
            }
            ICommand cmd;
            switch (args[1])
            {
                case "list":
                    var serviceNames = Services.OrderBy(s => s.ServiceType.Name).Select(s => s.ServiceType.Name);
                    cmd = new ListCommand(db, args.GetElementOrDefault(2), args.GetElementOrDefault(3), userManager, roleManager, serviceNames);
                    break;
            }
        }

        private static ServiceProvider GetServiceProvider()
        {
            // Get configuration from user secrets if in dev mode.
            var builder = new ConfigurationBuilder();
            string env = Environment.GetEnvironmentVariable("ASPNETCORE_ENVIRONMENT");
            if (env == "Development")
            {
                builder.AddUserSecrets<Program>();
            } // TODO: else use env variables
            Configuration = builder.Build();

            var services = new ServiceCollection();
            var conn = Configuration.GetConnectionString("Ry");
            services.AddDbContext<ApplicationDbContext>(options =>
                options.UseNpgsql(conn)
            );
            services.AddIdentity<ApplicationUser, ApplicationRole>()
                .AddEntityFrameworkStores<ApplicationDbContext>();

            Services = services;
            return services.BuildServiceProvider();
        }
    }
}
