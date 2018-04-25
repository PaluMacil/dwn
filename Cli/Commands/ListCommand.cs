using System;
using Microsoft.EntityFrameworkCore;
using Dwn.Data.Database;
using Dwn.Data.Models.Identity;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Configuration;
using Microsoft.AspNetCore.Identity;
using System.Linq;
using System.Collections.Generic;

namespace Dwn.Cli.Commands
{
    public class ListCommand : ICommand
    {
        private ApplicationDbContext Db;
        private string PrimarySubject;
        private string SecondarySubject;
        private UserManager<ApplicationUser> UserManager;
        private RoleManager<IdentityRole> RoleManager;
        private static IEnumerable<string> ServiceNames;
        

        public ListCommand(
            ApplicationDbContext db,
            string primarySubject, 
            string secondarySubject, 
            UserManager<ApplicationUser> userManager,
            RoleManager<IdentityRole> roleManager,
            IEnumerable<string> serviceNames)
        {
            Db = db;
            PrimarySubject = primarySubject.ToLower();
            SecondarySubject = secondarySubject.ToLower();
            UserManager = userManager;
            RoleManager = roleManager;
            ServiceNames = serviceNames;
        }

        public void Execute()
        {
            switch (PrimarySubject)
            {
                case "users":
                    Users(SecondarySubject);
                    break;
                case "roles":
                    Roles(SecondarySubject);
                    break;
                case "injectable":
                    Injectable();
                    break;
            }
        }

        private void Injectable()
        {
            foreach (var s in ServiceNames)
            {
                Console.WriteLine(s);
            }
        }

        private void Users(string forRole)
        {
            IList<ApplicationUser> users;
            if (forRole == null)
            {
                users = Db.Users.ToList();
            }
            else
            {
                users = UserManager.GetUsersInRoleAsync(forRole).Result;
            }
            
            foreach (var u in users)
            {
                Console.WriteLine(u.Email);
            }
        }

        private void Roles(string forUser)
        {

        }
    }
}
